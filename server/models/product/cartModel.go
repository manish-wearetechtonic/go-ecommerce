package product

import (
	user "github.com/thisismanishrajput/go-ecommerce/server/models"
)

type Cart struct {
	User     *user.User
	Products []struct {
		Product  *Product
		Quantity int `json:"quantity"`
	} `json:"cart"`
}
