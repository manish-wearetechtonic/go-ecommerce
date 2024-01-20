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
	"github.com/thisismanishrajput/go-ecommerce/server/models/product"
)

var brandCollection *mongo.Collection = database.OpenCollection(database.Client, "brand")

func AddBrand() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var brand product.Brand

		// Parse JSON request body into the brand struct
		if err := c.BindJSON(&brand); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate the brand struct
		validationError := validate.Struct(brand)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}

		// Check if the brand already exists
		count, err := brandCollection.CountDocuments(ctx, bson.M{"brand_name": brand.Brand_name})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while checking for existing brand"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Brand already exists"})
			return
		}

		// If the brand doesn't exist, insert it into the database
		brand.ID = primitive.NewObjectID()

		_, insertErr := brandCollection.InsertOne(ctx, brand)
		if insertErr != nil {
			log.Panic(insertErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while adding brand"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Brand added successfully", "brand": brand})
	}
}

func GetAllBrands() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// Find all brands in the database
		cursor, err := brandCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching brands"})
			return
		}
		defer cursor.Close(ctx)

		// Iterate over the cursor and collect brands
		var brands []product.Brand
		for cursor.Next(ctx) {
			var brand product.Brand
			if err := cursor.Decode(&brand); err != nil {
				log.Panic(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding brand"})
				return
			}
			brands = append(brands, brand)
		}

		if err := cursor.Err(); err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while iterating over brands"})
			return
		}

		// Respond with the list of brands
		c.JSON(http.StatusOK, gin.H{"success": 200, "message": "Given all brand", "brands": brands})
	}
}
