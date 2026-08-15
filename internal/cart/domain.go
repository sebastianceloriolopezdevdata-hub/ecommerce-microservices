package cart

import (
	"context"
	"fmt"

	"github.com/sebastianceloriolopez/ecommerce-microservices/cart-service/internal/product"
)

type CartItem struct {
	ProductID string  `json:"product_id"`
	Name      string  `json:"name,omitempty"`
	Price     float64 `json:"price,omitempty"`
	Category  string  `json:"category,omitempty"`
	Quantity  int     `json:"quantity"`
	TotalPrice float64 `json:"total_price,omitempty"`
}

type Cart struct {
	UserID string     `json:"user_id"`
	Items  []CartItem `json:"items"`
}

type CartStore interface {
	Get(ctx context.Context, userID string) (Cart, error)

	AddItem(
		ctx context.Context,
		userID string,
		item CartItem,
	) (Cart, error)

	UpdateItem(
		ctx context.Context,
		userID string,
		productID string,
		quantity int,
	) (Cart, error)

	RemoveItem(
		ctx context.Context,
		userID string,
		productID string,
	) (Cart, error)

	Clear(ctx context.Context, userID string) error
}

type productClient interface {
	GetProduct(productID string) (*product.Product, error)
}

type Service struct {
	store         CartStore
	productClient productClient
}

func NewService(store CartStore, productClient productClient) *Service {
	return &Service{
		store:         store,
		productClient: productClient,
	}
}

func (s *Service) GetCart(ctx context.Context, userID string) (Cart, error) {
	cart, err := s.store.Get(ctx, userID)
	if err != nil {
		return Cart{}, err
	}

	if s.productClient == nil {
		return cart, nil
	}

	validItems := make([]CartItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		productData, err := s.productClient.GetProduct(item.ProductID)
		if err != nil {
			continue
		}

		item.Name = productData.Name
		item.Price = productData.Price
		item.Category = productData.CategoryName()
		item.TotalPrice = productData.Price * float64(item.Quantity)
		validItems = append(validItems, item)
	}

	cart.Items = validItems
	return cart, nil
}

func (s *Service) AddItem(
	ctx context.Context,
	userID string,
	item CartItem,
) (Cart, error) {
	if s.productClient != nil {
		productData, err := s.productClient.GetProduct(item.ProductID)
		if err != nil {
			return Cart{}, fmt.Errorf("invalid product: %w", err)
		}

		item.Name = productData.Name
		item.Price = productData.Price
		item.Category = productData.CategoryName()
		item.TotalPrice = productData.Price * float64(item.Quantity)
	}

	return s.store.AddItem(ctx, userID, item)
}

func (s *Service) UpdateItem(
	ctx context.Context,
	userID string,
	productID string,
	quantity int,
) (Cart, error) {
	return s.store.UpdateItem(ctx, userID, productID, quantity)
}

func (s *Service) RemoveItem(
	ctx context.Context,
	userID string,
	productID string,
) (Cart, error) {
	return s.store.RemoveItem(ctx, userID, productID)
}

func (s *Service) ClearCart(ctx context.Context, userID string) error {
	return s.store.Clear(ctx, userID)
}