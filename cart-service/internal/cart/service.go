
package cart

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

type addItemRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type updateItemRequest struct {
	Quantity int `json:"quantity"`
}

const defaultCartID = "default"

// GET /api/cart
func (h *Handler) GetCart(
	w http.ResponseWriter,
	r *http.Request,
) {

	cart, err := h.svc.GetCart(
		r.Context(),
		defaultCartID,
	)

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"Error retrieving cart",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		cart,
	)
}

// POST /api/cart/items
func (h *Handler) AddItem(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req addItemRequest

	if err := json.NewDecoder(
		r.Body,
	).Decode(&req); err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)

		return
	}

	if req.ProductID == "" {

		writeError(
			w,
			http.StatusBadRequest,
			"Product ID is required",
		)

		return
	}

	if req.Quantity <= 0 {

		writeError(
			w,
			http.StatusBadRequest,
			"Quantity must be greater than zero",
		)

		return
	}

	cart, err := h.svc.AddItem(
		r.Context(),
		defaultCartID,
		CartItem{
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		},
	)

	if err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		cart,
	)
}

// PUT /api/cart/items/{productID}
func (h *Handler) UpdateItem(
	w http.ResponseWriter,
	r *http.Request,
) {

	productID := r.PathValue("productID")

	var req updateItemRequest

	if err := json.NewDecoder(
		r.Body,
	).Decode(&req); err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)

		return
	}

	if req.Quantity <= 0 {

		writeError(
			w,
			http.StatusBadRequest,
			"Quantity must be greater than zero",
		)

		return
	}

	cart, err := h.svc.UpdateItem(
		r.Context(),
		defaultCartID,
		productID,
		req.Quantity,
	)

	if err != nil {

		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		cart,
	)
}

// DELETE /api/cart/items/{productID}
func (h *Handler) RemoveItem(
	w http.ResponseWriter,
	r *http.Request,
) {

	productID := r.PathValue("productID")

	cart, err := h.svc.RemoveItem(
		r.Context(),
		defaultCartID,
		productID,
	)

	if err != nil {

		writeError(
			w,
			http.StatusInternalServerError,
			"Error removing product",
		)

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		cart,
	)
}

// DELETE /api/cart
func (h *Handler) ClearCart(
	w http.ResponseWriter,
	r *http.Request,
) {

	if err := h.svc.ClearCart(
		r.Context(),
		defaultCartID,
	); err != nil {

		writeError(
			w,
			http.StatusInternalServerError,
			"Error clearing cart",
		)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data interface{},
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func writeError(
	w http.ResponseWriter,
	status int,
	message string,
) {

	writeJSON(
		w,
		status,
		map[string]string{
			"message": message,
		},
	)
}