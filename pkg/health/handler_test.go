package health_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/health"
	"github.com/stretchr/testify/assert"
)

func TestHealth_OK(t *testing.T) {
	h := health.Handler{
		Service: &serviceMock{},
	}

	req, err := http.NewRequest("GET", "/health", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.Health(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestHealth_ServiceError(t *testing.T) {
	h := health.Handler{
		Service: &serviceMock{
			shouldFail: true,
		},
	}

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.Health(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

// mocks

type serviceMock struct {
	shouldFail     bool
	serviceStatus  bool
	externalStatus bool
	cacheStatus    bool
}

func (s *serviceMock) HealthCheck(ctx context.Context) (service, external, cache bool, err error) {
	if s.shouldFail {
		return s.serviceStatus, s.externalStatus, s.cacheStatus, fmt.Errorf("service asked to fail")
	}
	return s.serviceStatus, s.externalStatus, s.cacheStatus, nil
}
