package cart_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/logger"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/cart"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/item"
)

func TestCreateCartOK(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.CreateCart(context.TODO())

	if err != nil {
		t.Fatalf("Service not Expected to fail")
	}
}

func TestCreateCartCacheFail(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{
			shouldSetFail: true,
		},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.CreateCart(context.TODO())

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestGetCartOK(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.GetCart(context.TODO(), "testCartID")

	if err != nil {
		t.Fatalf("Service not Expected to fail")
	}
}

func TestGetCartCacheFail(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{
			shouldGetFail: true,
		},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.GetCart(context.TODO(), "testCartID")

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestGetCartExternalFail(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: true,
		})

	_, err := svc.GetCart(context.TODO(), "testCartID")

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestGetAvailableItemsOK(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.GetAvailableItems(context.TODO())

	if err != nil {
		t.Fatalf("Service not Expected to fail")
	}
}

func TestGetAvailableItemsExternalFailure(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: true,
		})

	_, err := svc.GetAvailableItems(context.TODO())

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestGetItemOK(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.GetItem(context.TODO(), "someItem")

	if err != nil {
		t.Fatalf("Service not Expected to fail")
	}
}

func TestGetItemExternalFailure(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: true,
		})

	_, err := svc.GetItem(context.TODO(), "someItem")

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestAddItemToCartOK(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.AddItemToCart(context.TODO(), "someCart", "someItem", 1)

	if err != nil {
		t.Fatalf("Service not Expected to fail")
	}
}

func TestAddItemToCartFailItemAlreadyAdded(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.AddItemToCart(context.TODO(), "someCart", "1-simple-Item", 1)

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestAddItemToCartCacheFailureGet(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{
			shouldGetFail: true,
		},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.AddItemToCart(context.TODO(), "someCart", "someItem", 1)

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestAddItemToCartCacheFailureSet(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{
			shouldSetFail: true,
		},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.AddItemToCart(context.TODO(), "someCart", "someItem", 1)

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestAddItemToCartExternalFailure(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: true,
		})

	_, err := svc.AddItemToCart(context.TODO(), "someCart", "someItem", 1)

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestModifyItemInCartOK(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.ModifyItemInCart(context.TODO(), "someCart", "1-simple-Item", 1)

	if err != nil {
		t.Fatalf("Service not Expected to fail")
	}
}

func TestModifyItemInCartItemNotFound(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.ModifyItemInCart(context.TODO(), "someCart", "SomeItem", 1)

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestModifyItemInCartCacheGetFailure(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{
			shouldGetFail: true,
		},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.ModifyItemInCart(context.TODO(), "someCart", "1-simple-Item", 1)

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestModifyItemInCartCacheSetFailure(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{
			shouldSetFail: true,
		},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.ModifyItemInCart(context.TODO(), "someCart", "1-simple-Item", 1)

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestModifyItemInCartExternalFailure(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: true,
		})

	_, err := svc.ModifyItemInCart(context.TODO(), "someCart", "1-simple-Item", 1)

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestDeleteItemInCartOK(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.DeleteItemInCart(context.TODO(), "someCart", "1-simple-Item")

	if err != nil {
		t.Fatalf("Service not Expected to fail")
	}
}

func TestDeleteItemInCartItemNotFound(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.DeleteItemInCart(context.TODO(), "someCart", "SomeItem")

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestDeleteItemInCartCacheGetFailure(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{
			shouldGetFail: true,
		},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.DeleteItemInCart(context.TODO(), "someCart", "1-simple-Item")

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestDeleteItemInCartCacheSetFailure(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{
			shouldSetFail: true,
		},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.DeleteItemInCart(context.TODO(), "someCart", "1-simple-Item")

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestDeleteItemInCartExternalFailure(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: true,
		})

	_, err := svc.DeleteItemInCart(context.TODO(), "someCart", "1-simple-Item")

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestDeleteAllItemsInCartOK(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.DeleteAllItemsInCart(context.TODO(), "someCart")

	if err != nil {
		t.Fatalf("Service not Expected to fail")
	}
}

func TestDeleteAllItemsInCartCacheGetFailure(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{
			shouldGetFail: true,
		},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.DeleteAllItemsInCart(context.TODO(), "someCart")

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestDeleteAllItemsInCartCacheSetFailure(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{
			shouldSetFail: true,
		},
		&externalMock{
			shouldFail: false,
		})

	_, err := svc.DeleteAllItemsInCart(context.TODO(), "someCart")

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

func TestDeleteCartOK(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{},
		&externalMock{
			shouldFail: false,
		})

	err := svc.DeleteCart(context.TODO(), "someCart")

	if err != nil {
		t.Fatalf("Service not Expected to fail")
	}
}

func TestDeleteCartCacheFailure(t *testing.T) {
	svc := cart.NewCartService("unit-testing",
		logger.NewLogger(false, "cart service unit testing"),
		&cacheMock{
			shouldDelFail: true,
		},
		&externalMock{
			shouldFail: false,
		})

	err := svc.DeleteCart(context.TODO(), "someCart")

	if err == nil {
		t.Fatalf("Service Expected to fail")
	}
}

//*************************Mocks********************

//******** Cache Mock

type cacheMock struct {
	shouldSetFail   bool
	shouldGetFail   bool
	shouldDelFail   bool
	shouldAliveFail bool
}

func (c *cacheMock) Set(ctx context.Context, key string, value interface{}) error {
	if c.shouldSetFail {
		return fmt.Errorf("Mock was asked to fail")
	}
	return nil
}
func (c *cacheMock) Get(ctx context.Context, key string, here interface{}) error {
	if c.shouldGetFail {
		return fmt.Errorf("Mock was asked to fail")
	}
	m := here.(*cart.Cart)

	m.Items = []item.Item{
		{
			ID: "1-simple-Item",
		},
		{
			ID: "2-simple-Item",
		},
	}
	return nil
}
func (c *cacheMock) Del(ctx context.Context, key string) error {
	if c.shouldDelFail {
		return fmt.Errorf("Mock was asked to fail")
	}

	return nil
}
func (c *cacheMock) Alive(ctx context.Context) bool {
	return !c.shouldAliveFail
}

//External Service Mock
type externalMock struct {
	shouldFail bool
}

func (e *externalMock) Health(ctx context.Context) error {
	if e.shouldFail {
		return fmt.Errorf("External API Mock was asked to fail")
	}
	return nil
}

func (e *externalMock) GetItem(ctx context.Context, id string) (item.Item, error) {
	if e.shouldFail {
		return item.Item{}, fmt.Errorf("External Mock was asked to Fail")
	}
	return item.Item{}, nil
}
func (e *externalMock) GetAllItems(ctx context.Context) ([]item.Item, error) {
	if e.shouldFail {
		return []item.Item{}, fmt.Errorf("External Mock was asked to Fail")
	}

	return []item.Item{}, nil
}
