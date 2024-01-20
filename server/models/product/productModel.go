package product

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents the product model in Go
type Product struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	Price           float64            `json:"price"`
	ProductCategory struct {
		CategoryID primitive.ObjectID `json:"category"`
		Name       string             `json:"name"`
	} `json:"productCategory"`
	ProductBrand struct {
		BrandID primitive.ObjectID `json:"brand"`
		Name    string             `json:"name"`
	} `json:"productBrand"`
	ImageUrls     []string `json:"imageUrls"`
	StockQuantity int      `json:"stockQuantity"`
	Ratings       []struct {
		UserID primitive.ObjectID `json:"user"`
		Rating float64            `json:"rating"`
		Review string             `json:"review"`
		Date   primitive.DateTime `json:"date"`
	} `json:"ratings"`
}

// Helper function to create a pointer to a string
func String(s string) *string {
	return &s
}
