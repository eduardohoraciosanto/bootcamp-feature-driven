package cart

import "github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/item"

type Cart struct {
	ID    string
	Items []item.Item
}

type TransportCart struct {
	ID    string               `json:"id"`
	Items []item.TransportItem `json:"items"`
}

type CartResponse struct {
	Cart TransportCart `json:"cart"`
}

func CartModelToTransportModel(cart Cart) TransportCart {
	vmItems := []item.TransportItem{}

	for _, i := range cart.Items {
		vmItems = append(vmItems, item.TransportItem{
			ID:       i.ID,
			Name:     i.Name,
			Quantity: i.Quantity,
			Price:    i.Price,
		})
	}

	return TransportCart{
		ID:    cart.ID,
		Items: vmItems,
	}
}

type AddItemToCartRequest struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type ModifyItemQuantityRequest struct {
	Quantity int `json:"quantity"`
}
