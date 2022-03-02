package errors

const (
	CartNotFoundCode           = "err_cart_not_found"
	ItemNotFoundCode           = "err_item_not_found"
	ItemNotFoundOnProviderCode = "err_provider_item_not_found"
	ItemAlreadyInCartCode      = "err_item_already_in_cart"
	ExternalApiErrorCode       = "err_external_api_error"
	CacheErrorCode             = "err_cache"
)

type ServiceError struct {
	Code string
}

func (s ServiceError) Error() string {
	return s.Code
}
