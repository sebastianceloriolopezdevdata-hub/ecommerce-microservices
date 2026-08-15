package cart

import (
	"context"
	"errors"
	"testing"

	"github.com/sebastianceloriolopez/ecommerce-microservices/cart-service/internal/product"
)

type fakeStore struct {
	called bool
	cart   Cart
}

func (f *fakeStore) Get(ctx context.Context, userID string) (Cart, error) {
	if f.cart.UserID == "" {
		return Cart{UserID: userID, Items: []CartItem{}}, nil
	}
	return f.cart, nil
}

func (f *fakeStore) AddItem(ctx context.Context, userID string, item CartItem) (Cart, error) {
	f.called = true
	return Cart{UserID: userID, Items: []CartItem{item}}, nil
}

func (f *fakeStore) UpdateItem(ctx context.Context, userID, productID string, quantity int) (Cart, error) {
	return Cart{}, nil
}

func (f *fakeStore) RemoveItem(ctx context.Context, userID, productID string) (Cart, error) {
	return Cart{}, nil
}

func (f *fakeStore) Clear(ctx context.Context, userID string) error {
	return nil
}

type fakeProductClient struct {
	err    error
	product *product.Product
}

func (f *fakeProductClient) GetProduct(productID string) (*product.Product, error) {
	if productID == "perro" {
		return nil, errors.New("product not found")
	}
	if f.err != nil {
		return nil, f.err
	}
	return f.product, nil
}

func TestAddItem_RejectsInvalidProductID(t *testing.T) {
	store := &fakeStore{}
	client := &fakeProductClient{err: errors.New("product not found")}
	svc := NewService(store, client)

	_, err := svc.AddItem(context.Background(), "default", CartItem{ProductID: "perro", Quantity: 2})
	if err == nil {
		t.Fatal("expected error for invalid product id")
	}

	if store.called {
		t.Fatal("store should not be called when product is invalid")
	}
}

func TestAddItem_StoresProductMetadata(t *testing.T) {
	store := &fakeStore{}
	client := &fakeProductClient{product: &product.Product{
		ID:         "123",
		Name:       "Laptop",
		Price:      999.99,
		CategoryID: map[string]interface{}{"name": "Electronics"},
	}}
	svc := NewService(store, client)

	cart, err := svc.AddItem(context.Background(), "default", CartItem{ProductID: "123", Quantity: 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(cart.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(cart.Items))
	}

	item := cart.Items[0]
	if item.Name != "Laptop" {
		t.Fatalf("expected name Laptop, got %q", item.Name)
	}
	if item.Price != 999.99 {
		t.Fatalf("expected price 999.99, got %v", item.Price)
	}
	if item.Category != "Electronics" {
		t.Fatalf("expected category Electronics, got %q", item.Category)
	}
	if item.TotalPrice != 1999.98 {
		t.Fatalf("expected total price 1999.98 for 2 items, got %v", item.TotalPrice)
	}
}

func TestGetCart_FiltersInvalidProductsAndHydratesMetadata(t *testing.T) {
	store := &fakeStore{cart: Cart{UserID: "default", Items: []CartItem{
		{ProductID: "123", Quantity: 2},
		{ProductID: "perro", Quantity: 1},
	}}}
	client := &fakeProductClient{product: &product.Product{
		ID:         "123",
		Name:       "Laptop",
		Price:      999.99,
		CategoryID: map[string]interface{}{"name": "Electronics"},
	}}
	svc := NewService(store, client)

	cart, err := svc.GetCart(context.Background(), "default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cart.Items) != 1 {
		t.Fatalf("expected 1 valid item, got %d items", len(cart.Items))
	}
	if cart.Items[0].Name != "Laptop" {
		t.Fatalf("expected hydrated name Laptop, got %q", cart.Items[0].Name)
	}
	if cart.Items[0].TotalPrice != 1999.98 {
		t.Fatalf("expected total price 1999.98, got %v", cart.Items[0].TotalPrice)
	}
}
