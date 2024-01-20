package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/thisismanishrajput/go-ecommerce/server/database"
	"github.com/thisismanishrajput/go-ecommerce/server/models"
	"github.com/thisismanishrajput/go-ecommerce/server/models/product"
)

var cartCollection *mongo.Collection = database.OpenCollection(database.Client, "cart")
var productCollection *mongo.Collection = database.OpenCollection(database.Client, "product")

func AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var request struct {
			ProductID string `json:"product_id" binding:"required"`
			Quantity  int    `json:"quantity" binding:"required"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the request fields
		validationError := validate.Struct(request)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}

		// Check if the user already has a cart
		userID := c.Param("userID")
		existingCart, err := getCartByUserID(ctx, userID)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for existing cart"})
			return
		}

		// If the user doesn't have a cart, create a new one
		if existingCart == nil {
			existingCart = &product.Cart{
				User: &models.User{User_id: userID},
			}
		}

		// Check if the product is already in the cart
		var found bool
		for i, item := range existingCart.Products {
			// Convert request.ProductID to primitive.ObjectID for comparison
			requestedProductID, err := primitive.ObjectIDFromHex(request.ProductID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
				return
			}

			if item.Product.ID == requestedProductID {
				existingCart.Products[i].Quantity += request.Quantity
				found = true
				break
			}
		}

		// If the product is not in the cart, add it
		if !found {
			// Fetch product details based on the product_id (you need to implement this)
			productDetails, err := getProductDetails(request.ProductID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching product details"})
				return
			}

			// Create a new cart item
			cartItem := product.CartItem{
				Product:  productDetails,
				Quantity: request.Quantity,
			}

			existingCart.Products = append(existingCart.Products, cartItem)
		}

		// Update the cart in the database
		err = updateCart(ctx, existingCart)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while updating cart"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully", "cart": existingCart})
	}
}

// getProductDetails fetches the product details based on the product ID (you need to implement this)
func getProductDetails(ctx context.Context, productID string) (*product.Product, error) {
	// Implement logic to fetch product details from your database or external service
	// Example: Fetch product details from the "products" collection in MongoDB
	var productDetails product.Product
	err := productCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&productDetails)
	if err != nil {
		return nil, err
	}

	return &productDetails, nil
}

// getCartByUserID retrieves the user's cart from the database
func getCartByUserID(ctx context.Context, userID string) (*product.Cart, error) {
	var existingCart product.Cart

	err := cartCollection.FindOne(ctx, bson.M{"user.user_id": userID}).Decode(&existingCart)
	if err == mongo.ErrNoDocuments {
		return nil, nil // User does not have a cart
	} else if err != nil {
		return nil, err
	}

	return &existingCart, nil
}

// updateCart updates the user's cart in the database
func updateCart(ctx context.Context, existingCart *product.Cart) error {
	_, err := cartCollection.UpdateOne(
		ctx,
		bson.M{"user.user_id": existingCart.User.User_id},
		bson.M{"$set": bson.M{"products": existingCart.Products}},
	)
	return err
}
