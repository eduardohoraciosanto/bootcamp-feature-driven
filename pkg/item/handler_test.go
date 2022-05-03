package item_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/item"
	"github.com/stretchr/testify/assert"
)

func TestGetAllItems_OK(t *testing.T) {
	h := item.Handler{
		Service: &mockedService{},
	}

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.GetAllItems(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetAllItems_Error(t *testing.T) {
	h := item.Handler{
		Service: &mockedService{
			shouldFail: true,
		},
	}

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.GetAllItems(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestGetItem_OK(t *testing.T) {
	h := item.Handler{
		Service: &mockedService{},
	}

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.GetItem(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetItem_Error(t *testing.T) {
	h := item.Handler{
		Service: &mockedService{
			shouldFail: true,
		},
	}

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.GetItem(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

// Mocks

type mockedService struct {
	shouldFail bool
}

func (m *mockedService) Health(ctx context.Context) error {
	if m.shouldFail {
		return fmt.Errorf("mock was asked to fail")
	}
	return nil
}
func (m *mockedService) GetItem(ctx context.Context, id string) (item.Item, error) {
	if m.shouldFail {
		return item.Item{}, fmt.Errorf("mock was asked to fail")
	}
	return item.Item{
		ID:       "someID",
		Name:     "someName",
		Quantity: 1,
		Price:    12.34,
	}, nil
}
func (m *mockedService) GetAllItems(ctx context.Context) ([]item.Item, error) {
	if m.shouldFail {
		return []item.Item{}, fmt.Errorf("mock was asked to fail")
	}
	return []item.Item{
		{
			ID:       "someID",
			Name:     "someName",
			Quantity: 1,
			Price:    12.34,
		},
	}, nil
}
