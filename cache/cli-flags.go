package cache

import (
	"github.com/sky-mirror/boot"
	"github.com/urfave/cli/v2"
)

type config struct {
	Addr             string
	Password         string
	MasterName       string
	SentinelPassword string
}

var defaultConfig config

var _ boot.Beforer = &config{}

func init() {
	boot.Register(&defaultConfig)
}

// CliFlags returns cli flags to setup cache package.
func (cfg *config) CliFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, &cli.StringFlag{
		Name:        "redis-addr",
		Value:       "redis:6379",
		EnvVars:     []string{"REDIS_ADDR"},
		Destination: &cfg.Addr,
	})

	flags = append(flags, &cli.StringFlag{
		Name:        "redis-password",
		EnvVars:     []string{"REDIS_PASSWORD"},
		Destination: &cfg.Password,
	})

	flags = append(flags, &cli.StringFlag{
		Name:        "redis-master-name",
		Destination: &cfg.MasterName,
		Usage:       "enables sentinel mode and specify master",
	})

	flags = append(flags, &cli.StringFlag{
		Name:        "redis-sentinel-password",
		EnvVars:     []string{"REDIS_SENTINEL_PASSWORD"},
		Destination: &cfg.SentinelPassword,
	})

	return flags
}

// Before inits.
func (cfg *config) Before(c *cli.Context) error {
	if len(cfg.MasterName) > 0 {
		InitializeSentinel(
			cfg.Addr,
			cfg.Password,
			cfg.MasterName,
			cfg.SentinelPassword,
		)
	} else {
		Initialize(cfg.Addr, cfg.Password)
	}

	return nil
}
