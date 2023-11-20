package helloer

import (
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"

	"go-monorepo/env"
	"go-monorepo/logging"
)

var myName = env.String(env.PodName)

// Service handles hello request and generates response.
type Service struct {
	db            *gorm.DB
	reqCounterVec *prometheus.CounterVec
}

// NewService creates service with db to write log.
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
		reqCounterVec: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: "hello",
			Subsystem: "service",
			Name:      "request_count_total",
			Help:      "Counter to request",
		}, []string{"request_type"}),
	}
}

func (s *Service) logResponse(typ, name, resp string) {
	if s.db == nil {
		return
	}

	if err := s.db.Create(&RequestLog{
		Type:     typ,
		Request:  name,
		Response: resp,
	}).Error; err != nil {
		logging.Get().Error(err)
	}
}

// SayHello generates response.
func (s *Service) SayHello(name string) (resp string, err error) {
	defer func() {
		s.logResponse("SayHello", name, resp)
	}()

	if s.reqCounterVec != nil {
		s.reqCounterVec.With(prometheus.Labels{"request_type": "SayHello"}).
			Inc()
	}

	switch name {
	case "deer":
		resp = "slap"
	case "DEADBEEF":
		err = errors.Wrapf(ErrBadName, "name(%v)", name)
	default:
		resp = "Hello, " + name
	}

	resp += " I'm " + *myName
	return
}
