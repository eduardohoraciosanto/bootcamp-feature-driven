package health

import (
	"context"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/cache"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/logger"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/item"
)

//Service is the interface for the health
type Service interface {
	HealthCheck(ctx context.Context) (service, external, cache bool, err error)
}

type svc struct {
	log   logger.Logger
	cache cache.Cache
	item  item.Service
}

//NewService gives a new Service
func NewService(c cache.Cache, item item.Service, log logger.Logger) Service {
	return &svc{
		log:   log,
		cache: c,
		item:  item,
	}
}

//HealthCheck returns the status of the API and it's components
func (s *svc) HealthCheck(ctx context.Context) (service, external, cache bool, err error) {
	s.log.Info(ctx, "Performing Health Check")

	external = true
	if err := s.item.Health(ctx); err != nil {
		external = false
	}
	return true, external, s.cache.Alive(ctx), nil
}
