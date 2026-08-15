package cart

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{
		rdb: rdb,
	}
}

func cartKey(userID string) string {
	return fmt.Sprintf("cart:%s", userID)
}

func (s *RedisStore) Get(
	ctx context.Context,
	userID string,
) (Cart, error) {

	key := cartKey(userID)

	data, err := s.rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return Cart{
			UserID: userID,
			Items:  []CartItem{},
		}, nil
	}

	if err != nil {
		return Cart{}, err
	}

	var cart Cart

	if err := json.Unmarshal(
		[]byte(data),
		&cart,
	); err != nil {
		return Cart{}, err
	}

	return cart, nil
}

func (s *RedisStore) AddItem(
	ctx context.Context,
	userID string,
	item CartItem,
) (Cart, error) {

	cart, err := s.Get(ctx, userID)

	if err != nil {
		return Cart{}, err
	}

	for i := range cart.Items {

		if cart.Items[i].ProductID == item.ProductID {
			cart.Items[i].Quantity += item.Quantity
			cart.Items[i].TotalPrice = cart.Items[i].Price * float64(cart.Items[i].Quantity)
			return s.save(ctx, cart)
		}
	}

	cart.Items = append(
		cart.Items,
		item,
	)

	return s.save(ctx, cart)
}

func (s *RedisStore) UpdateItem(
	ctx context.Context,
	userID string,
	productID string,
	quantity int,
) (Cart, error) {

	cart, err := s.Get(ctx, userID)

	if err != nil {
		return Cart{}, err
	}

	for i := range cart.Items {

		if cart.Items[i].ProductID == productID {
			cart.Items[i].Quantity = quantity
			cart.Items[i].TotalPrice = cart.Items[i].Price * float64(quantity)
			return s.save(ctx, cart)
		}
	}

	return Cart{}, fmt.Errorf("product not found in cart")
}

func (s *RedisStore) RemoveItem(
	ctx context.Context,
	userID string,
	productID string,
) (Cart, error) {

	cart, err := s.Get(ctx, userID)

	if err != nil {
		return Cart{}, err
	}

	items := make([]CartItem, 0, len(cart.Items))

	for _, item := range cart.Items {

		if item.ProductID != productID {
			items = append(items, item)
		}
	}

	cart.Items = items

	return s.save(ctx, cart)
}

func (s *RedisStore) Clear(
	ctx context.Context,
	userID string,
) error {

	return s.rdb.Del(
		ctx,
		cartKey(userID),
	).Err()
}

func (s *RedisStore) save(
	ctx context.Context,
	cart Cart,
) (Cart, error) {

	data, err := json.Marshal(cart)

	if err != nil {
		return Cart{}, err
	}

	err = s.rdb.Set(
		ctx,
		cartKey(cart.UserID),
		data,
		0,
	).Err()

	if err != nil {
		return Cart{}, err
	}

	return cart, nil
}