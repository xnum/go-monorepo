package logging

import (
	"errors"
	"os"
	"path"

	"github.com/sky-mirror/boot"
	"github.com/urfave/cli/v2"
)

var defaultCfg config

func init() {
	boot.Register(&defaultCfg)
}

type fileConfig struct {
	enabled bool
	dir     string
	name    string
}

type fluentConfig struct {
	enabled bool
	host    string
	port    int
	tag     string
}

type config struct {
	file   fileConfig
	fluent fluentConfig
}

var _ boot.Beforer = &config{}
var _ boot.Afterer = &config{}

func (cfg *config) CliFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, &cli.BoolFlag{
		Name:        "log-enable-file",
		EnvVars:     []string{"LOG_ENABLE_FILE"},
		Destination: &cfg.file.enabled,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        "log-file-name",
		EnvVars:     []string{"LOG_FILE_NAME"},
		Usage:       "filename prefix of log file",
		Value:       path.Base(os.Args[0]),
		Destination: &cfg.file.name,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        "log-file-dir",
		EnvVars:     []string{"LOG_FILE_DIR"},
		Usage:       "path of log file",
		Value:       os.TempDir(),
		Destination: &cfg.file.dir,
	})

	flags = append(flags, &cli.BoolFlag{
		Name:        "log-enable-fluent",
		EnvVars:     []string{"LOG_ENABLE_FLUENT"},
		Destination: &cfg.fluent.enabled,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        "log-fluent-host",
		EnvVars:     []string{"LOG_FLUENT_HOST"},
		Destination: &cfg.fluent.host,
	})
	flags = append(flags, &cli.IntFlag{
		Name:        "log-fluent-port",
		EnvVars:     []string{"LOG_FLUENT_PORT"},
		Destination: &cfg.fluent.port,
	})
	flags = append(flags, &cli.StringFlag{
		Name:        "log-fluent-tag",
		EnvVars:     []string{"LOG_FLUENT_TAG"},
		Destination: &cfg.fluent.tag,
	})

	return flags
}

func (cfg *config) Before(c *cli.Context) error {
	if cfg.file.enabled {
		if len(cfg.file.dir) == 0 {
			return errors.New("log-file-dir must be set")
		}
		if len(cfg.file.name) == 0 {
			return errors.New("log-file-name must be set")
		}
	}

	if cfg.fluent.enabled {
		if len(cfg.fluent.host) == 0 {
			return errors.New("log-fluent-host must be set")
		}
		if len(cfg.fluent.tag) == 0 {
			return errors.New("log-fluent-tag must be set")
		}
		if cfg.fluent.port == 0 {
			return errors.New("log-fluent-port must be set")
		}
	}

	Initialize()

	return nil
}

func (cfg *config) After() {
	Finalize()
}

// TestingInitialize inits package and writes log to temp dir.
func TestingInitialize() {
	Initialize()
}

// TestingFinalize removes closes file and remove it.
func TestingFinalize() {
	Finalize()
	// TODO: remove log
}
