package item

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/errors"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/logger"
)

const (
	healthEndpoint   = "https://bootcamp-products.getsandbox.com/health"
	articlesEndpoint = "https://bootcamp-products.getsandbox.com/products"

	healthStatusOK = "OK"
)

type Service interface {
	Health(ctx context.Context) error
	GetItem(ctx context.Context, id string) (Item, error)
	GetAllItems(ctx context.Context) ([]Item, error)
}

type externalService struct {
	client ItemClient
	logger logger.Logger
}

type ItemClient interface {
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (*http.Response, error)
}

func NewExternalService(logger logger.Logger, client ItemClient) Service {

	return &externalService{
		logger: logger,
		client: client,
	}
}

func (e *externalService) Health(ctx context.Context) error {
	e.logger.Info(ctx, "Calling External API Health")
	outCtx := context.WithValue(context.Background(), "X-Correlation-ID", ctx.Value("correlation_id"))
	req, err := http.NewRequestWithContext(outCtx, http.MethodGet, healthEndpoint, nil)
	if err != nil {
		e.logger.WithError(err).Error(ctx, "Error creating request to external provider")
		return err
	}
	res, err := e.client.Do(req)
	if err != nil {
		e.logger.WithError(err).Error(ctx, "Error Calling External API Health")
		return err
	}
	eHealth := ExternalHealthResponse{}

	err = json.NewDecoder(res.Body).Decode(&eHealth)
	if err != nil {
		e.logger.WithError(err).Error(ctx, "Error Decoding External API Health")
		return err
	}

	if eHealth.Data.Status != healthStatusOK {
		e.logger.WithField("external_api_status", eHealth.Data.Status).Error(ctx, "External API Not Healthy")
		return fmt.Errorf("external API not Healthy - Status: %s", eHealth.Data.Status)
	}
	return nil
}
func (e *externalService) GetItem(ctx context.Context, id string) (Item, error) {
	log := e.logger.WithField("item_id", id)

	log.Info(ctx, "Getting single item from provider")
	outCtx := context.WithValue(context.Background(), "X-Correlation-ID", ctx.Value("correlation_id"))

	req, err := http.NewRequestWithContext(outCtx, http.MethodGet, articlesEndpoint+"/"+id, nil)
	if err != nil {
		log.WithError(err).Error(ctx, "Error creating request to external provider")
		return Item{}, err
	}
	res, err := e.client.Do(req)
	if err != nil {
		log.WithError(err).Error(ctx, "Error getting item from external provider")
		return Item{}, err
	}
	if res.StatusCode == http.StatusNotFound {
		log.Error(ctx, "Item from external provider not found")
		return Item{}, errors.ServiceError{Code: errors.ItemNotFoundOnProviderCode}
	}
	eItem := ExternalGetItemResponse{}

	err = json.NewDecoder(res.Body).Decode(&eItem)
	if err != nil {
		log.WithError(err).Error(ctx, "Unable to parse provider response")
		return Item{}, err
	}
	price, err := strconv.ParseFloat(eItem.Data.Price, 32)
	if err != nil {
		log.WithError(err).Error(ctx, "Unable to parse price into a floating point number")
		return Item{}, err
	}
	log.Info(ctx, "Item fetched successfully")
	mItem := Item{
		ID:    eItem.Data.ID,
		Name:  eItem.Data.Name,
		Price: float32(price),
	}

	return mItem, nil
}
func (e *externalService) GetAllItems(ctx context.Context) ([]Item, error) {
	e.logger.Info(ctx, "Getting all items from provider")
	outCtx := context.WithValue(context.Background(), "X-Correlation-ID", ctx.Value("correlation_id"))
	req, err := http.NewRequestWithContext(outCtx, http.MethodGet, articlesEndpoint, nil)
	if err != nil {
		e.logger.WithError(err).Error(ctx, "Error creating request to external provider")
		return []Item{}, err
	}
	res, err := e.client.Do(req)
	if err != nil {
		e.logger.WithError(err).Error(ctx, "Error getting all items from external provider")
		return []Item{}, err
	}
	eItems := ExternalGetAllItemsResponse{}

	err = json.NewDecoder(res.Body).Decode(&eItems)
	if err != nil {
		e.logger.WithError(err).Error(ctx, "Unable to parse provider response")
		return []Item{}, err
	}

	mItems := []Item{}
	for _, eItem := range eItems.Data {
		price, err := strconv.ParseFloat(eItem.Price, 32)
		if err != nil {
			e.logger.WithError(err).Error(ctx, "Unable to parse price into a floating point number")
			return []Item{}, err
		}
		mItems = append(mItems, Item{
			ID:    eItem.ID,
			Name:  eItem.Name,
			Price: float32(price),
		})
	}

	return mItems, nil
}
