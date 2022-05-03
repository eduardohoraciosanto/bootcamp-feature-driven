package health

import (
	"context"
	"fmt"
	"testing"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/logger"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/item"
)

func TestHealthCheck(t *testing.T) {
	service := NewService(
		&cacheMocked{cacheShouldFail: false},
		&externalAPIMocked{externalAPIShouldFail: false},
		logger.NewLogger("health svc unit test", false),
	)
	s, e, d, err := service.HealthCheck(context.TODO())
	if s != true || e != true || d != true || err != nil {
		t.Errorf("Unexpected values from method: service %t, db %t, error %s", s, d, err)
	}
}

func TestHealthCheck_CacheFail(t *testing.T) {
	service := NewService(
		&cacheMocked{cacheShouldFail: true},
		&externalAPIMocked{externalAPIShouldFail: false},
		logger.NewLogger("health svc unit test", false),
	)

	s, e, d, err := service.HealthCheck(context.TODO())
	if s != true || e != true || d != false || err != nil {
		t.Errorf("Unexpected values from method: service %t, db %t, error %s", s, d, err)
	}
}

func TestHealthCheck_ExternalFail(t *testing.T) {
	service := NewService(
		&cacheMocked{cacheShouldFail: false},
		&externalAPIMocked{externalAPIShouldFail: true},
		logger.NewLogger("health svc unit test", false),
	)

	s, e, d, err := service.HealthCheck(context.TODO())
	if s != true || e != false || d != true || err != nil {
		t.Errorf("Unexpected values from method: service %t, db %t, error %s", s, d, err)
	}
}

//Cache Mocked

type cacheMocked struct {
	cacheShouldFail bool
}

func (c *cacheMocked) Set(ctx context.Context, key string, value interface{}) error {
	if c.cacheShouldFail {
		return fmt.Errorf("Mock Cache Asked to Fail")
	}
	return nil
}
func (c *cacheMocked) Get(ctx context.Context, key string, here interface{}) error {
	if c.cacheShouldFail {
		return fmt.Errorf("Mock Cache Asked to Fail")
	}
	return nil
}
func (c *cacheMocked) Del(ctx context.Context, key string) error {
	if c.cacheShouldFail {
		return fmt.Errorf("Mock Cache Asked to Fail")
	}
	return nil
}
func (c *cacheMocked) Alive(ctx context.Context) bool {
	if c.cacheShouldFail {
		return false
	}
	return true
}

type externalAPIMocked struct {
	externalAPIShouldFail bool
}

func (e *externalAPIMocked) Health(ctx context.Context) error {
	if e.externalAPIShouldFail {
		return fmt.Errorf("External API Mock was asked to fail")
	}
	return nil
}

func (e *externalAPIMocked) GetItem(ctx context.Context, id string) (item.Item, error) {
	if e.externalAPIShouldFail {
		return item.Item{}, fmt.Errorf("External API Mock was asked to fail")
	}
	return item.Item{
		ID:    "mockedItem",
		Name:  "Mocked Item",
		Price: 999.99,
	}, nil
}
func (e *externalAPIMocked) GetAllItems(ctx context.Context) ([]item.Item, error) {
	if e.externalAPIShouldFail {
		return []item.Item{}, fmt.Errorf("External API Mock was asked to fail")
	}
	return []item.Item{
		{
			ID:    "mockedItem1",
			Name:  "Mocked Item 1",
			Price: 999.99,
		},
		{
			ID:    "mockedItem2",
			Name:  "Mocked Item 2",
			Price: 999.99,
		},
	}, nil
}
