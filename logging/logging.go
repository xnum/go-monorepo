package logging

import (
	"log"
	"os"
	"path/filepath"

	"github.com/op/go-logging"

	"go-monorepo/pkg/fluent"
)

var (
	logger *logging.Logger
	file   *os.File
)

func levelToSyslogLevel(lvl logging.Level) int {
	switch lvl {
	case logging.CRITICAL:
		return 2
	case logging.ERROR:
		return 3
	case logging.WARNING:
		return 4
	case logging.NOTICE:
		return 5
	case logging.INFO:
		return 6
	case logging.DEBUG:
		return 7
	default:
		return 7
	}
}

func levelToLower(lvl logging.Level) string {
	switch lvl {
	case logging.CRITICAL:
		return "critical"
	case logging.ERROR:
		return "error"
	case logging.WARNING:
		return "warning"
	case logging.NOTICE:
		return "notice"
	case logging.INFO:
		return "info"
	case logging.DEBUG:
		return "debug"
	default:
		return "trace"
	}
}

type fluentBackend struct {
	client *fluent.Client
}

func (b *fluentBackend) Log(lvl logging.Level, calldepth int,
	record *logging.Record) error {

	data := map[string]any{}
	data["level"] = levelToSyslogLevel(record.Level)
	data["severity"] = levelToLower(record.Level)
	data["message"] = record.Message()

	b.client.Push(data)
	return nil
}

// Initialize inits singleton.
func Initialize() {
	logger = logging.MustGetLogger("app")
	format := logging.MustStringFormatter(
		`%{color}%{time:2006-01-02 15:04:05.000} %{shortfile} %{shortfunc}() %{level:.5s} %{color:reset} %{message}`,
	)

	errbackend := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Formatter := logging.NewBackendFormatter(errbackend, format)
	errBackendLeveled := logging.AddModuleLevel(backend1Formatter)

	var backends []logging.Backend
	backends = append(backends, errBackendLeveled)

	if defaultCfg.fluent.enabled {
		fluentBackend := &fluentBackend{
			client: fluent.NewClient(
				defaultCfg.fluent.host,
				defaultCfg.fluent.port,
				defaultCfg.fluent.tag,
			),
		}
		backends = append(backends, fluentBackend)
	}

	if defaultCfg.file.enabled {
		var err error
		file, err = os.OpenFile(
			filepath.Join(defaultCfg.file.dir, defaultCfg.file.name+".log"),
			os.O_RDWR|os.O_CREATE,
			0644,
		)
		if err != nil {
			log.Panicf(
				"OpenFile((%v,%v), ...): %v",
				defaultCfg.file.dir,
				defaultCfg.file.name,
				err,
			)
		}

		filebackend := logging.NewLogBackend(file, "", 0)
		backend2Formatter := logging.NewBackendFormatter(filebackend, format)
		fileBackendLeveled := logging.AddModuleLevel(backend2Formatter)

		backends = append(backends, fileBackendLeveled)
	}

	logging.SetBackend(backends...)
}

// Finalize closes singleton.
func Finalize() {
	if file != nil {
		file.Close()
	}
}

// Get gets logger.
func Get() *logging.Logger {
	if logger == nil {
		log.Panicln("uninit logging")
	}
	return logger
}
