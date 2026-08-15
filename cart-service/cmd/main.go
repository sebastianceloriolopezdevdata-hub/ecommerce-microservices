package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"

	"github.com/sebastianceloriolopez/ecommerce-microservices/cart-service/internal/cart"
	"github.com/sebastianceloriolopez/ecommerce-microservices/cart-service/internal/product"
)

func main() {

	redisURL := os.Getenv("REDIS_URL")

	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	productServiceURL := os.Getenv(
		"PRODUCT_SERVICE_URL",
	)

	if productServiceURL == "" {
		productServiceURL = "http://localhost:3000"
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "3001"
	}

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf(
			"Redis connection failed: %v",
			err,
		)
	}

	log.Println("Redis connection successful")

	// Product Service client
	productClient := product.NewClient(
		productServiceURL,
	)

	// Redis repository
	store := cart.NewRedisStore(rdb)

	// Service
	service := cart.NewService(
		store,
		productClient,
	)

	// HTTP handler
	handler := cart.NewHandler(service)

	// Routes
	mux := http.NewServeMux()

	cart.RegisterRoutes(
		mux,
		handler,
	)

	// Health check
	mux.HandleFunc(
		"GET /",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			writeHealth(w)
		},
	)

	log.Printf(
		"Cart Service running on port %s",
		port,
	)

	// Wrap mux with middleware (logging first, then CORS)
	httpHandler := cart.LoggingMiddleware(
		cart.CORSMiddleware(mux),
	)

	if err := http.ListenAndServe(
		":"+port,
		httpHandler,
	); err != nil {

		log.Fatal(err)
	}
}

func writeHealth(w http.ResponseWriter) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(`{
		"ok": true,
		"message": "Cart Service API running"
	}`))
}