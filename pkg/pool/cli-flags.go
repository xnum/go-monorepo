package pool

import (
	"strings"
	"time"

	"github.com/sky-mirror/boot"
	"github.com/urfave/cli/v2"
)

const (
	defaultInitCap     = 5
	defaultMaxCap      = 10
	defaultDialTimeout = 5 * time.Second
	defaultIdleTimeout = 60 * time.Second
)

// Config defines config to pool.
type Config struct {
	InitCap     int
	MaxCap      int
	DialTimeout time.Duration
	IdleTimeout time.Duration

	prefix string
}

// NewConfig creates config and register it.
func NewConfig(prefix string) *Config {
	if len(prefix) > 0 && !strings.HasSuffix(prefix, "-") {
		prefix += "-"
	}

	cfg := &Config{prefix: prefix}
	boot.Register(cfg)
	return cfg
}

// CliFlags returns cli flags to setup package.
func (cfg *Config) CliFlags() []cli.Flag {
	var flags []cli.Flag
	// define common value in env var or specific value in args.
	flags = append(flags, &cli.IntFlag{
		Name:        cfg.prefix + "pool-init-cap",
		EnvVars:     []string{"POOL_INIT_CAP"},
		Destination: &cfg.InitCap,
		Value:       defaultInitCap,
	})
	flags = append(flags, &cli.IntFlag{
		Name:        cfg.prefix + "pool-max-cap",
		EnvVars:     []string{"POOL_MAX_CAP"},
		Destination: &cfg.MaxCap,
		Value:       defaultMaxCap,
	})
	flags = append(flags, &cli.DurationFlag{
		Name:        cfg.prefix + "pool-dial-timeout",
		EnvVars:     []string{"POOL_DIAL_TIMEOUT"},
		Destination: &cfg.DialTimeout,
		Value:       defaultDialTimeout,
	})
	flags = append(flags, &cli.DurationFlag{
		Name:        cfg.prefix + "pool-idle-timeout",
		EnvVars:     []string{"POOL_IDLE_TIMEOUT"},
		Destination: &cfg.IdleTimeout,
		Value:       defaultIdleTimeout,
	})

	return flags
}
