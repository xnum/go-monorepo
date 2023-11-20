package env

import (
	"fmt"
	"os"

	"github.com/sky-mirror/boot"
	"github.com/urfave/cli/v2"
)

var keys map[string]*string

// Hook add hooks to cliflag.
type Hook struct{}

// CliFlags returns nothing.
func (h *Hook) CliFlags() []cli.Flag {
	return nil
}

// Before lookups env.
func (h *Hook) Before(*cli.Context) error {
	// lookup envs.
	for k, v := range keys {
		s, ok := os.LookupEnv(k)
		if !ok {
			return fmt.Errorf("required env var(%v) not found", k)
		}

		*v = s
	}

	return nil
}

func init() {
	keys = make(map[string]*string)

	boot.Register(&Hook{})
}

// well known keys.
const (
	PodName string = "POD_NAME"
	PodIP   string = "POD_IP"
)

// String registers key which is required and extracts from env to returned
// string pointer.
func String(key string) *string {
	if _, ok := keys[key]; !ok {
		keys[key] = new(string)
	}

	return keys[key]
}
