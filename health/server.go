package health

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go-monorepo/logging"
)

const (
	collectPeriod = 125 * time.Millisecond
	pollPeriod    = 3 * time.Second
)

var infos sync.Map

var tmpInfo struct {
	sync.Mutex
	alive bool
	ready bool
	vars  map[string]any
}

// StartCollector starts collector.
func StartCollector() {
	for {
		var (
			alive = true
			ready = true
			vars  = make(map[string]any)
		)

		infos.Range(func(key, value any) bool {
			info := value.(Info)

			age := time.Since(info.lastTime)
			expired := info.interval > 0 && age > info.interval
			if info.ProbeType != ProbeNone {
				if expired || info.Status != statRunning {
					switch info.ProbeType {
					case ProbeAlive:
						alive = false
					case ProbeReady:
						ready = false
					}
				}
				if info.Status == statExited || info.Status == statInit {
					log.Printf("set alive probe to false due to task(%v) is %v",
						info.taskName, info.Status)
					alive = false
				}
			}

			info.Age = age.String()
			info.Expired = expired
			vars[info.taskName] = info

			return true
		})

		if len(vars) == 0 {
			alive = true
			ready = true
		}

		vars["CollectedAt"] = time.Now()

		tmpInfo.Lock()
		tmpInfo.alive = alive
		tmpInfo.ready = ready
		tmpInfo.vars = vars
		tmpInfo.Unlock()

		time.Sleep(collectPeriod)
	}
}

func gatherInfos() (alive, ready bool, vars map[string]any) {
	tmpInfo.Lock()
	defer tmpInfo.Unlock()

	return tmpInfo.alive, tmpInfo.ready, tmpInfo.vars
}

func aliveHandler(ctx *gin.Context) {
	defer ctx.Abort()
	isAlive, _, _ := gatherInfos()

	if isAlive {
		ctx.Status(http.StatusOK)
	} else {
		ctx.Status(http.StatusServiceUnavailable)
	}
}

func readyHandler(ctx *gin.Context) {
	defer ctx.Abort()
	_, isReady, _ := gatherInfos()

	if isReady {
		ctx.Status(http.StatusOK)
	} else {
		ctx.Status(http.StatusServiceUnavailable)
	}
}

func varHandler(ctx *gin.Context) {
	defer ctx.Abort()

	code := http.StatusOK
	_, _, vars := gatherInfos()

	ctx.JSON(code, vars)
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// RegisterToGinEngine registers health check endpoint to an existing engine.
func RegisterToGinEngine(engine *gin.Engine) {
	engine.GET("/alive", aliveHandler)
	engine.GET("/ready", readyHandler)
	engine.GET("/vars", varHandler)
	engine.GET("/metrics", prometheusHandler())
}

var engine *gin.Engine

// Engine returns engine.
func Engine() *gin.Engine {
	return engine
}

func init() {
	engine = gin.New()
	engine.RedirectTrailingSlash = true
	RegisterToGinEngine(engine)
}

// StartServer starts health server and blocks.
func StartServer() {
	go StartCollector()

	srv := &http.Server{
		Addr:    defaultConfig.HealthAddr,
		Handler: engine,
	}

	logging.Get().Info("starts serving health server at", srv.Addr)
	if err := srv.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}
