package product

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Product struct {
	ID         string      `json:"_id"`
	Name       string      `json:"name"`
	Price      float64     `json:"price"`
	Stock      int         `json:"stock"`
	CategoryID interface{} `json:"categoryId"`
	Attributes map[string]interface{} `json:"attributes"`
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) GetProduct(productID string) (*Product, error) {

	url := fmt.Sprintf(
		"%s/api/products/%s",
		c.BaseURL,
		productID,
	)

	log.Printf("[Product Client] GET %s", url)

	resp, err := c.HTTPClient.Get(url)

	if err != nil {
		log.Printf("[Product Client] Error: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	log.Printf("[Product Client] Status: %d", resp.StatusCode)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("product not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"product service returned status %d",
			resp.StatusCode,
		)
	}

	var product Product

	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		log.Printf("[Product Client] Decode error: %v", err)
		return nil, err
	}

	log.Printf("[Product Client] Success: %+v", product)

	return &product, nil
}

func (p *Product) CategoryName() string {
	categoryName := ""

	switch category := p.CategoryID.(type) {
	case string:
		categoryName = category
	case map[string]interface{}:
		if name, ok := category["name"].(string); ok {
			categoryName = name
		}
	}

	return categoryName
}
