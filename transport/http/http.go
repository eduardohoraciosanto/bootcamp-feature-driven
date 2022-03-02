package transport

import (
	"context"
	"net/http"

	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/cart"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/health"
	"github.com/eduardohoraciosanto/bootcamp-feature-driven/pkg/item"
	"github.com/google/uuid"

	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

func NewHTTPRouter(hsvc health.Service, csvc cart.Service, isvc item.Service) *muxtrace.Router {

	hc := health.Handler{
		Service: hsvc,
	}

	cc := cart.Handler{
		Service: csvc,
	}

	ic := item.Handler{
		Service: isvc,
	}

	r := muxtrace.NewRouter()
	r.Use(correlationIDMiddleware)

	r.HandleFunc("/health", hc.Health).Methods(http.MethodGet)

	//Cart Endpoints
	r.HandleFunc("/cart", cc.CreateCart).Methods(http.MethodPost)
	r.HandleFunc("/cart/{cart_id}", cc.GetCart).Methods(http.MethodGet)
	r.HandleFunc("/cart/{cart_id}", cc.DeleteCart).Methods(http.MethodDelete)

	//Item Operations on Cart
	r.HandleFunc("/cart/{cart_id}/item", cc.AddItem).Methods(http.MethodPost)
	r.HandleFunc("/cart/{cart_id}/item/{item_id:[0-9]+}", cc.UpdateQuantity).Methods(http.MethodPut)
	r.HandleFunc("/cart/{cart_id}/item/all", cc.RemoveAllItems).Methods(http.MethodDelete)
	r.HandleFunc("/cart/{cart_id}/item/{item_id:[0-9]+}", cc.RemoveItem).Methods(http.MethodDelete)

	//Items Endpoints
	r.HandleFunc("/items/available", ic.GetAllItems).Methods(http.MethodGet)
	r.HandleFunc("/items/{item_id}", ic.GetItem).Methods(http.MethodGet)

	r.PathPrefix("/swagger").Handler(http.StripPrefix("/swagger", http.FileServer(http.Dir("./swagger"))))
	return r
}

func correlationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.Header.Get("X-Correlation-Id")
		if id == "" {
			// generate new version 4 uuid
			id = uuid.New().String()
		}
		// set the id to the request context
		ctx = context.WithValue(ctx, "correlation_id", id)
		r = r.WithContext(ctx)

		// set the response header
		w.Header().Set("X-Correlation-Id", id)
		next.ServeHTTP(w, r)
	})
}
