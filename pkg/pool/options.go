package pool

import (
	"errors"
	"math/rand"
	"runtime/debug"
	"time"

	"go-monorepo/logging"
)

var (
	errClosed   = errors.New("pool is closed")
	errInvalid  = errors.New("invalid config")
	errRejected = errors.New("connection is nil. rejecting")
	errTargets  = errors.New("targets server is empty")
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// Options pool options
type Options struct {
	//InitTargets init targets
	InitTargets []string
	// init connection
	InitCap int
	// max connections
	MaxCap      int
	DialTimeout time.Duration
	IdleTimeout time.Duration
}

// DefaultOptions returns a new newOptions instance with sane defaults.
func DefaultOptions(addrs ...string) *Options {
	o := &Options{}
	o.InitTargets = addrs
	o.InitCap = defaultInitCap
	o.MaxCap = defaultMaxCap
	o.DialTimeout = defaultDialTimeout
	o.IdleTimeout = defaultIdleTimeout
	return o
}

// NewOptions creates options from config.
func NewOptions(cfg *Config, addrs []string) *Options {
	return &Options{
		InitTargets: addrs,
		InitCap:     cfg.InitCap,
		MaxCap:      cfg.MaxCap,
		DialTimeout: cfg.DialTimeout,
		IdleTimeout: cfg.IdleTimeout,
	}
}

// validate checks a Config instance.
func (o *Options) validate() error {
	if len(o.InitTargets) == 0 ||
		o.InitCap <= 0 ||
		o.MaxCap <= 0 ||
		o.InitCap > o.MaxCap ||
		o.DialTimeout == 0 {
		debug.PrintStack()
		logging.Get().Info("options", o)
		return errInvalid
	}
	return nil
}

// nextTarget next target implement load balance
func (o *Options) nextTarget() string {
	tlen := len(o.InitTargets)

	//rand server
	return o.InitTargets[rand.Int()%tlen]
}
