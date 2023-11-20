package hello

import (
	"strings"
	"time"

	"github.com/sky-mirror/boot"
	"github.com/urfave/cli/v2"

	"go-monorepo/pkg/pool"
)

// ClientConfig defines parameter to initialize hello client.
type ClientConfig struct {
	prefix string

	HelloServiceEndpoint string
	CallTimeout          time.Duration

	poolConfig *pool.Config
}

// NewClientConfig creates client config and register it.
func NewClientConfig(prefix string) *ClientConfig {
	if len(prefix) > 0 && !strings.HasSuffix(prefix, "-") {
		prefix += "-"
	}

	cfg := &ClientConfig{prefix: prefix}
	boot.Register(cfg)

	cfg.poolConfig = pool.NewConfig(prefix + "hello")
	return cfg
}

// CliFlags returns cli flags to setup package.
func (cfg *ClientConfig) CliFlags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, &cli.StringFlag{
		Name:        cfg.prefix + "hello-service-endpoint",
		EnvVars:     []string{"QUOTE_SERVICE_ENDPOINT"},
		Destination: &cfg.HelloServiceEndpoint,
		Required:    true,
	})
	flags = append(flags, &cli.DurationFlag{
		Name:        cfg.prefix + "hello-pool-call-timeout",
		EnvVars:     []string{"POOL_CALL_TIMEOUT"},
		Destination: &cfg.CallTimeout,
		Value:       3 * time.Second,
	})

	return flags
}
