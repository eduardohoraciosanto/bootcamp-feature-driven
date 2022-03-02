package cart

import (
	"context"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/cache"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/errors"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/logger"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/item"
	"github.com/google/uuid"
)

type Service interface {
	CreateCart(ctx context.Context) (Cart, error)
	GetCart(ctx context.Context, cartID string) (Cart, error)
	GetAvailableItems(ctx context.Context) ([]item.Item, error)
	GetItem(ctx context.Context, id string) (item.Item, error)
	AddItemToCart(ctx context.Context, cartID, itemID string, quantity int) (Cart, error)
	ModifyItemInCart(ctx context.Context, cartID, itemID string, newQuantity int) (Cart, error)
	DeleteItemInCart(ctx context.Context, cartID, itemID string) (Cart, error)
	DeleteAllItemsInCart(ctx context.Context, cartID string) (Cart, error)
	DeleteCart(ctx context.Context, cartID string) error
}

type service struct {
	//dependencies of the service
	logger          logger.Logger
	version         string
	cache           cache.Cache
	externalService item.Service
}

func NewCartService(version string, logger logger.Logger, cache cache.Cache, externalService item.Service) Service {
	return &service{
		logger:          logger,
		version:         version,
		cache:           cache,
		externalService: externalService,
	}
}

func (s *service) CreateCart(ctx context.Context) (Cart, error) {
	cartID := uuid.New().String()

	log := s.logger.WithField("cart_id", cartID)

	cart := Cart{
		ID: cartID,
	}
	log.Info(ctx, "Creating new cart")
	if err := s.cache.Set(ctx, cartID, cart); err != nil {
		log.WithError(err).Error(ctx, "Unable to save new cart in DB")
		return Cart{}, errors.ServiceError{
			Code: errors.CacheErrorCode,
		}
	}

	return cart, nil
}

func (s *service) GetCart(ctx context.Context, cartID string) (Cart, error) {
	log := s.logger.WithField("cart_id", cartID)

	cart := Cart{}
	log.Info(ctx, "Getting cart from DB")
	err := s.cache.Get(ctx, cartID, &cart)
	if err != nil {
		log.WithError(err).Error(ctx, "Unable to save new cart in DB")
		return Cart{}, errors.ServiceError{Code: errors.CartNotFoundCode}
	}
	log.WithField("cart_id", cartID).Info(ctx, "Populating items info from provider")
	err = s.fetchItemsForCart(ctx, &cart)
	if err != nil {
		log.WithError(err).Error(ctx, "Error fetching items for cart")
		return Cart{}, errors.ServiceError{Code: errors.ExternalApiErrorCode}
	}

	return cart, nil
}

func (s *service) GetAvailableItems(ctx context.Context) ([]item.Item, error) {
	log := s.logger

	log.Info(ctx, "Fetching all items from provider")
	items, err := s.externalService.GetAllItems(ctx)
	if err != nil {
		log.WithError(err).Error(ctx, "Error fetching items from provider")
		return []item.Item{}, err
	}
	return items, nil
}

func (s *service) GetItem(ctx context.Context, id string) (item.Item, error) {
	log := s.logger.WithField("item_id", id)

	log.Info(ctx, "Fetching item from provider")
	i, err := s.externalService.GetItem(ctx, id)
	if err != nil {
		log.WithError(err).Error(ctx, "Error fetching item from provider")
		return item.Item{}, err
	}
	return i, nil
}

func (s *service) AddItemToCart(ctx context.Context, cartID, itemID string, quantity int) (Cart, error) {
	log := s.logger.
		WithField("cart_id", cartID).
		WithField("item_id", itemID).
		WithField("quantity", quantity)

	log.Info(ctx, "Adding Item to Cart")
	cart := Cart{}
	log.Info(ctx, "Getting Cart from DB")
	err := s.cache.Get(ctx, cartID, &cart)
	if err != nil {
		log.WithError(err).Error(ctx, "Unable to get Cart from DB")
		return Cart{}, errors.ServiceError{Code: errors.CartNotFoundCode}
	}

	log.Info(ctx, "Adding item to Cart")
	for _, item := range cart.Items {
		if item.ID == itemID {
			log.WithError(err).Error(ctx, "Item Already in Cart")
			return Cart{}, errors.ServiceError{Code: errors.ItemAlreadyInCartCode}
		}
	}

	cart.Items = append(cart.Items, item.Item{
		ID:       itemID,
		Quantity: quantity,
	})

	log.Info(ctx, "Saving Cart to DB")
	if err := s.cache.Set(ctx, cartID, cart); err != nil {
		log.WithError(err).Error(ctx, "Unable to save Cart in DB")
		return Cart{}, err
	}

	log.Info(ctx, "Getting Cart Item details from provider")
	err = s.fetchItemsForCart(ctx, &cart)
	if err != nil {
		log.WithError(err).Error(ctx, "Unable to get data from the provider")
		return Cart{}, errors.ServiceError{Code: errors.ExternalApiErrorCode}
	}

	return cart, nil
}
func (s *service) ModifyItemInCart(ctx context.Context, cartID, itemID string, newQuantity int) (Cart, error) {
	log := s.logger.
		WithField("cart_id", cartID).
		WithField("item_id", itemID).
		WithField("new_quantity", newQuantity)

	log.Info(ctx, "Modifying item quantity in Cart")

	cart := Cart{}

	log.Info(ctx, "Getting Cart from DB")
	err := s.cache.Get(ctx, cartID, &cart)
	if err != nil {
		log.WithError(err).Error(ctx, "Unable to get Cart from DB")
		return Cart{}, errors.ServiceError{Code: errors.CartNotFoundCode}
	}

	log.Info(ctx, "Looking for Item in Cart")
	for idx, item := range cart.Items {
		if item.ID == itemID {
			log.Info(ctx, "Updating Item in Cart")
			cart.Items[idx].Quantity = newQuantity
			log.Info(ctx, "Saving Cart in DB")
			if err := s.cache.Set(ctx, cartID, cart); err != nil {
				log.WithError(err).Error(ctx, "Unable Saving Cart to DB")
				return Cart{}, err
			}
			log.Info(ctx, "Getting Cart Item details from provider")
			err = s.fetchItemsForCart(ctx, &cart)
			if err != nil {
				log.WithError(err).Error(ctx, "Unable to fetch item data from provider")
				return Cart{}, errors.ServiceError{Code: errors.ExternalApiErrorCode}
			}
			return cart, nil
		}
	}
	log.Error(ctx, "Unable to find Item inside Cart")
	return Cart{}, errors.ServiceError{Code: errors.ItemNotFoundCode}
}
func (s *service) DeleteItemInCart(ctx context.Context, cartID, itemID string) (Cart, error) {
	log := s.logger.
		WithField("cart_id", cartID).
		WithField("item_id", itemID)

	log.Info(ctx, "Deleting item from Cart")
	cart := Cart{}
	log.Info(ctx, "Getting Cart from DB")
	err := s.cache.Get(ctx, cartID, &cart)
	if err != nil {
		log.WithError(err).Error(ctx, "Unable to get Cart from DB")
		return Cart{}, errors.ServiceError{Code: errors.CartNotFoundCode}
	}

	log.Info(ctx, "Removing item from Cart data")
	for idx, item := range cart.Items {
		if item.ID == itemID {
			//we care about the order, so we perform to sub-slices

			cart.Items = append(cart.Items[:idx], cart.Items[idx+1:]...)

			log.Info(ctx, "Saving Cart in DB")
			if err := s.cache.Set(ctx, cartID, cart); err != nil {
				return Cart{}, err
			}
			log.Info(ctx, "Getting Cart Item details from provider")
			err = s.fetchItemsForCart(ctx, &cart)
			if err != nil {
				log.WithError(err).Error(ctx, "Unable to fetch item data from provider")
				return Cart{}, errors.ServiceError{Code: errors.ExternalApiErrorCode}
			}

			return cart, nil
		}
	}
	log.WithError(err).Error(ctx, "Unable to find Item inside Cart")
	return Cart{}, errors.ServiceError{Code: errors.ItemNotFoundCode}
}
func (s *service) DeleteAllItemsInCart(ctx context.Context, cartID string) (Cart, error) {
	log := s.logger.WithField("cart_id", cartID)
	log.Info(ctx, "Deleting all items in Cart")

	cart := Cart{}
	log.Info(ctx, "Getting Cart from DB")
	err := s.cache.Get(ctx, cartID, &cart)
	if err != nil {
		log.WithError(err).Error(ctx, "Unable to get Cart from DB")
		return Cart{}, errors.ServiceError{Code: errors.CartNotFoundCode}
	}

	log.Info(ctx, "Removing all items from Cart data")
	cart.Items = []item.Item{}
	log.Info(ctx, "Saving Cart in DB")
	if err := s.cache.Set(ctx, cartID, cart); err != nil {
		log.WithError(err).Error(ctx, "Unable to save Cart to DB")
		return Cart{}, err
	}

	return cart, nil
}
func (s *service) DeleteCart(ctx context.Context, cartID string) error {
	log := s.logger.WithField("cart_id", cartID)

	log.Info(ctx, "Deleting Cart entirely")
	err := s.cache.Del(ctx, cartID)
	if err != nil {
		log.WithError(err).Error(ctx, "Unable to delete Cart from DB")
		return errors.ServiceError{Code: errors.CartNotFoundCode}
	}
	return nil
}
func (s *service) fetchItemsForCart(ctx context.Context, cart *Cart) error {
	log := s.logger.WithField("cart_id", cart.ID)

	log.Info(ctx, "Fetching Cart's items from provider")
	//We fetch information from the external service to fill in Name and Price
	for idx, item := range cart.Items {
		extItem, err := s.externalService.GetItem(ctx, item.ID)
		if err != nil {
			log.WithField("item_id", item.ID).Error(ctx, "Unable to get item from provider")
			return err
		}
		cart.Items[idx].Price = extItem.Price
		cart.Items[idx].Name = extItem.Name
	}
	return nil
}
