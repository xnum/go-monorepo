package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"go-monorepo/appmodule/counter"
	"go-monorepo/logging"
)

// Server is a HTTP server.
type Server struct {
	service *counter.Service
}

// RegisterMiddleware registers middleware for all endpoints.
func (s *Server) RegisterMiddleware(r *gin.Engine) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	config := cors.Config{
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"HEAD",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	config.AllowAllOrigins = true

	r.Use(cors.New(config))
	r.Use(RateLimit())
}

// RegisterEndpoint installs api representation layer processing function.
func (s *Server) RegisterEndpoint(r *gin.Engine) {
	r.GET("count", func(ctx *gin.Context) {
		count := s.service.Query()

		ctx.JSON(200, gin.H{"count": count})
	})
}

// Start starts HTTP server.
func (s *Server) Start(ctx context.Context, apiAddr string) {
	gin.ForceConsoleColor()

	// setup gin.
	apiEngine := gin.New()
	apiEngine.RedirectTrailingSlash = true

	s.RegisterMiddleware(apiEngine)

	// setup service.
	s.service = &counter.Service{
		// can do DI here.
	}
	s.service.Start(ctx)

	// setup endpoint.
	s.RegisterEndpoint(apiEngine)

	srv := &http.Server{
		Addr:    apiAddr,
		Handler: apiEngine,
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown: ", err)
		}
	}()

	logging.Get().Info("starts serving...")
	if err := srv.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}
