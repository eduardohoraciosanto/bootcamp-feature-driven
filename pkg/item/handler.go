package item

import (
	"net/http"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/internal/response"
	"github.com/gorilla/mux"
)

type Handler struct {
	Service Service
}

//GetAllItems returns all items from the external API
func (c *Handler) GetAllItems(w http.ResponseWriter, r *http.Request) {

	items, err := c.Service.GetAllItems(r.Context())
	if err != nil {
		response.RespondWithError(w, err)
		return
	}
	vmItems := []TransportItem{}
	for _, item := range items {
		vmItem := TransportItem{
			ID:    item.ID,
			Name:  item.Name,
			Price: item.Price,
		}
		vmItems = append(vmItems, vmItem)
	}

	response.RespondWithData(w, http.StatusOK, vmItems)
}

//GetItem returns a particular item from the external API
func (c *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID := vars["item_id"]
	item, err := c.Service.GetItem(r.Context(), itemID)
	if err != nil {
		response.RespondWithError(w, err)
		return
	}
	vmItem := TransportItem{
		ID:    item.ID,
		Name:  item.Name,
		Price: item.Price,
	}

	response.RespondWithData(w, http.StatusOK, vmItem)
}
