package cart

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/response"
	"github.com/gorilla/mux"
)

type Handler struct {
	Service Service
}

//CreateCart creates a cart on the DB
func (c *Handler) CreateCart(w http.ResponseWriter, r *http.Request) {

	cart, err := c.Service.CreateCart(r.Context())
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	res := CartResponse{
		Cart: CartModelToTransportModel(cart),
	}
	response.RespondWithData(w, http.StatusOK, res)
}

//GetCart creates a cart on the DB
func (c *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID := vars["cart_id"]
	cart, err := c.Service.GetCart(r.Context(), cartID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}
	res := CartResponse{
		Cart: CartModelToTransportModel(cart),
	}
	response.RespondWithData(w, http.StatusOK, res)
}

//DeleteCart removes all items from the cart
func (c *Handler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID := vars["cart_id"]

	err := c.Service.DeleteCart(r.Context(), cartID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}
	response.RespondWithData(w, http.StatusAccepted, nil)
}

//AddItem Adds an item to a cart
func (c *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID := vars["cart_id"]
	vm := AddItemToCartRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&vm)
	if err != nil {
		log.Printf("Error decoding body: %v", err)
		response.RespondWithError(w, response.StandardBadBodyRequest)
		return
	}
	cart, err := c.Service.AddItemToCart(r.Context(), cartID, vm.ID, vm.Quantity)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	res := CartResponse{
		Cart: CartModelToTransportModel(cart),
	}
	response.RespondWithData(w, http.StatusOK, res)
}

//UpdateQuantity changes the amount of a single item in the cart
func (c *Handler) UpdateQuantity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID := vars["cart_id"]
	itemID := vars["item_id"]

	vm := ModifyItemQuantityRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&vm)
	if err != nil {
		log.Printf("Error decoding body: %v", err)
		response.RespondWithError(w, response.StandardBadBodyRequest)
		return
	}

	cart, err := c.Service.ModifyItemInCart(r.Context(), cartID, itemID, vm.Quantity)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	res := CartResponse{
		Cart: CartModelToTransportModel(cart),
	}
	response.RespondWithData(w, http.StatusOK, res)
}

//RemoveItem removes an item from the cart
func (c *Handler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID := vars["cart_id"]
	itemID := vars["item_id"]

	cart, err := c.Service.DeleteItemInCart(r.Context(), cartID, itemID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	res := CartResponse{
		Cart: CartModelToTransportModel(cart),
	}

	response.RespondWithData(w, http.StatusOK, res)
}

//RemoveAllItems removes all items from the cart
func (c *Handler) RemoveAllItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cartID := vars["cart_id"]

	cart, err := c.Service.DeleteAllItemsInCart(r.Context(), cartID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}

	res := CartResponse{
		Cart: CartModelToTransportModel(cart),
	}
	response.RespondWithData(w, http.StatusOK, res)
}
