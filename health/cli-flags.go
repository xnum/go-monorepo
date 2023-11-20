package health

import (
	"github.com/sky-mirror/boot"
	"github.com/urfave/cli/v2"
)

type config struct {
	HealthAddr string
}

var defaultConfig config

func init() {
	boot.Register(&defaultConfig)
}

// CliFlags returns cli flags to setup cache package.
func (cfg *config) CliFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, &cli.StringFlag{
		Name:        "health-addr",
		Value:       ":7001",
		EnvVars:     []string{"HEALTH_ADDR"},
		Destination: &cfg.HealthAddr,
	})

	return flags
}
