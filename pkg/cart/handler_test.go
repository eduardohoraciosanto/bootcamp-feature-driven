package cart_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/cart"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/item"
	"github.com/stretchr/testify/assert"
)

func TestCreateCart_OK(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{},
	}

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.CreateCart(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestCreateCart_Error(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{
			shouldFail: true,
		},
	}

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.CreateCart(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestGetCart_OK(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{},
	}

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.GetCart(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetCart_Error(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{
			shouldFail: true,
		},
	}

	req, err := http.NewRequest("GET", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.GetCart(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestDeleteCart_OK(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{},
	}

	req, err := http.NewRequest("DELETE", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.DeleteCart(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusAccepted, res.StatusCode)
}

func TestDeleteCart_Error(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{
			shouldFail: true,
		},
	}

	req, err := http.NewRequest("DELETE", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.DeleteCart(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestAddItemToCart_OK(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{},
	}

	payload := cart.AddItemToCartRequest{
		ID:       "someID",
		Quantity: 1,
	}
	pBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(pBytes))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.AddItem(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestAddItemToCart_BadPayload(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{},
	}

	pBytes := []byte("notJSON")

	req, err := http.NewRequest("POST", "/", bytes.NewReader(pBytes))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.AddItem(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestAddItemToCart_Error(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{
			shouldFail: true,
		},
	}

	payload := cart.AddItemToCartRequest{
		ID:       "someID",
		Quantity: 1,
	}
	pBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(pBytes))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.AddItem(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestUpdateQuantity_OK(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{},
	}

	payload := cart.ModifyItemQuantityRequest{
		Quantity: 1,
	}
	pBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(pBytes))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.UpdateQuantity(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestUpdateQuantity_BadPayload(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{},
	}

	pBytes := []byte("NotJSON")

	req, err := http.NewRequest("POST", "/", bytes.NewReader(pBytes))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.UpdateQuantity(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestUpdateQuantity_Error(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{
			shouldFail: true,
		},
	}

	payload := cart.ModifyItemQuantityRequest{
		Quantity: 1,
	}
	pBytes, err := json.Marshal(payload)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(pBytes))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.UpdateQuantity(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestRemoveItem_OK(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{},
	}

	req, err := http.NewRequest("DELETE", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.RemoveItem(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestRemoveItem_Error(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{
			shouldFail: true,
		},
	}

	req, err := http.NewRequest("DELETE", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.RemoveItem(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestRemoveAllItems_OK(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{},
	}

	req, err := http.NewRequest("DELETE", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.RemoveAllItems(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestRemoveAllItems_Error(t *testing.T) {
	h := cart.Handler{
		Service: &mockedService{
			shouldFail: true,
		},
	}

	req, err := http.NewRequest("DELETE", "/", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()

	h.RemoveAllItems(rr, req)

	res := rr.Result()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

// Mocks

type mockedService struct {
	shouldFail bool
}

func (m *mockedService) CreateCart(ctx context.Context) (cart.Cart, error) {
	if m.shouldFail {
		return cart.Cart{}, fmt.Errorf("mock was asked to fail")
	}

	return cart.Cart{}, nil
}
func (m *mockedService) GetCart(ctx context.Context, cartID string) (cart.Cart, error) {
	if m.shouldFail {
		return cart.Cart{}, fmt.Errorf("mock was asked to fail")
	}
	return cart.Cart{
		ID: cartID,
		Items: []item.Item{
			{
				ID:       "someItemID",
				Name:     "someItemName",
				Quantity: 2,
				Price:    12.34,
			},
		},
	}, nil
}
func (m *mockedService) GetAvailableItems(ctx context.Context) ([]item.Item, error) {
	if m.shouldFail {
		return []item.Item{}, fmt.Errorf("mock was asked to fail")
	}
	return []item.Item{}, nil
}
func (m *mockedService) GetItem(ctx context.Context, id string) (item.Item, error) {
	if m.shouldFail {
		return item.Item{}, fmt.Errorf("mock was asked to fail")
	}
	return item.Item{
		ID: id,
	}, nil
}
func (m *mockedService) AddItemToCart(ctx context.Context, cartID, itemID string, quantity int) (cart.Cart, error) {
	if m.shouldFail {
		return cart.Cart{}, fmt.Errorf("mock was asked to fail")
	}
	return cart.Cart{
		ID: cartID,
	}, nil
}
func (m *mockedService) ModifyItemInCart(ctx context.Context, cartID, itemID string, newQuantity int) (cart.Cart, error) {
	if m.shouldFail {
		return cart.Cart{}, fmt.Errorf("mock was asked to fail")
	}
	return cart.Cart{
		ID: cartID,
	}, nil
}
func (m *mockedService) DeleteItemInCart(ctx context.Context, cartID, itemID string) (cart.Cart, error) {
	if m.shouldFail {
		return cart.Cart{}, fmt.Errorf("mock was asked to fail")
	}
	return cart.Cart{
		ID: cartID,
	}, nil
}
func (m *mockedService) DeleteAllItemsInCart(ctx context.Context, cartID string) (cart.Cart, error) {
	if m.shouldFail {
		return cart.Cart{}, fmt.Errorf("mock was asked to fail")
	}
	return cart.Cart{
		ID: cartID,
	}, nil
}
func (m *mockedService) DeleteCart(ctx context.Context, cartID string) error {
	if m.shouldFail {
		return fmt.Errorf("mock was asked to fail")
	}

	return nil
}
