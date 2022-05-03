package item_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/logger"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/item"
)

func TestHealthOK(t *testing.T) {

	svc := item.NewExternalService(
		logger.NewLogger("item unit test", false),
		&itemClientMock{
			shouldFail: false,
			response: item.ExternalHealthResponse{
				Meta: item.ExternalMeta{
					Version: "testing",
				},
				Data: item.ExternalHealth{
					Status: "OK",
				},
			},
		},
	)

	err := svc.Health(context.TODO())
	if err != nil {
		t.Fatalf("Error was not expected")
	}
}

func TestHealthStatusIncorrect(t *testing.T) {

	svc := item.NewExternalService(
		logger.NewLogger("item unit test", false),
		&itemClientMock{
			shouldFail: false,
			response: item.ExternalHealthResponse{
				Meta: item.ExternalMeta{
					Version: "testing",
				},
				Data: item.ExternalHealth{
					Status: "Incorrect",
				},
			},
		},
	)

	err := svc.Health(context.TODO())
	if err == nil {
		t.Fatalf("Error was expected")
	}
}
func TestHealthError(t *testing.T) {

	svc := item.NewExternalService(
		logger.NewLogger("item unit test", false),
		&itemClientMock{
			shouldFail: true,
			response:   nil,
		},
	)

	err := svc.Health(context.TODO())
	if err == nil {
		t.Fatalf("Error was expected")
	}
}
func TestHealthDecodingError(t *testing.T) {

	svc := item.NewExternalService(
		logger.NewLogger("item unit test", false),
		&itemClientMock{
			shouldFail: false,
			response:   "notAJSON",
		},
	)

	err := svc.Health(context.TODO())
	if err == nil {
		t.Fatalf("Error was expected")
	}
}

func TestGetItem(t *testing.T) {

	svc := item.NewExternalService(
		logger.NewLogger("item unit test", false),
		&itemClientMock{
			shouldFail: false,
			response: item.ExternalGetItemResponse{
				Meta: item.ExternalMeta{
					Version: "testing",
				},
				Data: item.ExternalItem{
					ID:    "someItemID",
					Name:  "Some Item ID",
					Price: "12.34",
				},
			},
		},
	)

	_, err := svc.GetItem(context.TODO(), "someItemID")
	if err != nil {
		t.Fatalf("Error was not expected")
	}
}
func TestGetItemNotFound(t *testing.T) {

	svc := item.NewExternalService(
		logger.NewLogger("item unit test", false),
		&itemClientMock{
			shouldFail:         false,
			response:           nil,
			responseStatusCode: 404,
		},
	)

	_, err := svc.GetItem(context.TODO(), "someItemID")
	if err == nil {
		t.Fatalf("Error was not expected")
	}
}
func TestGetItemApiFailure(t *testing.T) {

	svc := item.NewExternalService(
		logger.NewLogger("item unit test", false),
		&itemClientMock{
			shouldFail: true,
			response:   nil,
		},
	)

	_, err := svc.GetItem(context.TODO(), "someItemID")
	if err == nil {
		t.Fatalf("Error was expected")
	}
}
func TestGetItemWrongResponse(t *testing.T) {

	svc := item.NewExternalService(
		logger.NewLogger("item unit test", false),
		&itemClientMock{
			shouldFail: false,
			response:   "WrongResponse",
		},
	)

	_, err := svc.GetItem(context.TODO(), "someItemID")
	if err == nil {
		t.Fatalf("Error was expected")
	}
}
func TestGetItemFloatParseError(t *testing.T) {

	svc := item.NewExternalService(
		logger.NewLogger("item unit test", false),
		&itemClientMock{
			shouldFail: false,
			response: item.ExternalGetItemResponse{
				Meta: item.ExternalMeta{
					Version: "testing",
				},
				Data: item.ExternalItem{
					ID:    "someItemID",
					Name:  "Some Item ID",
					Price: "notANumber",
				},
			},
		},
	)

	_, err := svc.GetItem(context.TODO(), "someItemID")
	if err == nil {
		t.Fatalf("Error was expected")
	}
}

func TestGetAllItems(t *testing.T) {

	svc := item.NewExternalService(
		logger.NewLogger("item unit test", false),
		&itemClientMock{
			shouldFail: false,
			response: item.ExternalGetAllItemsResponse{
				Meta: item.ExternalMeta{
					Version: "testing",
				},
				Data: []item.ExternalItem{
					{
						ID:    "someItemID",
						Name:  "Some Item ID",
						Price: "12.34",
					},
				},
			},
		},
	)

	_, err := svc.GetAllItems(context.TODO())
	if err != nil {
		t.Fatalf("Error was not expected")
	}
}
func TestGetAllItemsApiFailure(t *testing.T) {

	svc := item.NewExternalService(
		logger.NewLogger("item unit test", false),
		&itemClientMock{
			shouldFail: true,
			response:   nil,
		},
	)

	_, err := svc.GetAllItems(context.TODO())
	if err == nil {
		t.Fatalf("Error was expected")
	}
}
func TestGetAllItemsWrongResponse(t *testing.T) {

	svc := item.NewExternalService(
		logger.NewLogger("item unit test", false),
		&itemClientMock{
			shouldFail: false,
			response:   "WrongResponse",
		},
	)

	_, err := svc.GetAllItems(context.TODO())
	if err == nil {
		t.Fatalf("Error was expected")
	}
}
func TestGetAllItemsFloatParseError(t *testing.T) {

	svc := item.NewExternalService(
		logger.NewLogger("item unit test", false),
		&itemClientMock{
			shouldFail: false,
			response: item.ExternalGetAllItemsResponse{
				Meta: item.ExternalMeta{
					Version: "testing",
				},
				Data: []item.ExternalItem{
					{
						ID:    "someItemID",
						Name:  "Some Item ID",
						Price: "notANumber",
					},
				},
			},
		},
	)

	_, err := svc.GetAllItems(context.TODO())
	if err == nil {
		t.Fatalf("Error was expected")
	}
}

//*****ItemClientMock

type itemClientMock struct {
	response           interface{}
	responseStatusCode int
	shouldFail         bool
}

func (i *itemClientMock) Get(url string) (*http.Response, error) {
	if i.shouldFail {
		return nil, fmt.Errorf("Mock asked to fail")
	}
	b, _ := json.Marshal(i.response)
	resp := &http.Response{}
	resp.Body = ioutil.NopCloser(bytes.NewReader(b))
	if i.responseStatusCode != 0 {
		resp.StatusCode = i.responseStatusCode
	}
	return resp, nil
}

func (i *itemClientMock) Do(req *http.Request) (*http.Response, error) {
	if i.shouldFail {
		return nil, fmt.Errorf("Mock asked to fail")
	}
	b, _ := json.Marshal(i.response)
	resp := &http.Response{}
	resp.Body = ioutil.NopCloser(bytes.NewReader(b))
	if i.responseStatusCode != 0 {
		resp.StatusCode = i.responseStatusCode
	}
	return resp, nil
}
