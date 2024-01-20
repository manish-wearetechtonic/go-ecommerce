package product

import (
	user "github.com/thisismanishrajput/go-ecommerce/server/models"
)

// CartItem represents an item in the cart
type CartItem struct {
	Product  *Product
	Quantity int `json:"quantity"`
}

// Cart represents the user's shopping cart
type Cart struct {
	User     *user.User
	Products []CartItem `json:"cart"`
}
