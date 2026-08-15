package cart

import "net/http"

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
) {

	mux.HandleFunc(
		"GET /api/cart",
		handler.GetCart,
	)

	mux.HandleFunc(
		"POST /api/cart/items",
		handler.AddItem,
	)

	mux.HandleFunc(
		"PUT /api/cart/items/{productID}",
		handler.UpdateItem,
	)

	mux.HandleFunc(
		"DELETE /api/cart/items/{productID}",
		handler.RemoveItem,
	)

	mux.HandleFunc(
		"DELETE /api/cart",
		handler.ClearCart,
	)
}