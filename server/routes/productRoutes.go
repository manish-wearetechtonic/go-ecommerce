package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/thisismanishrajput/go-ecommerce/server/controllers"
)

func ProductRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/brand", controller.AddBrand())
	incomingRoutes.GET("/brand", controller.GetAllBrands())
}
