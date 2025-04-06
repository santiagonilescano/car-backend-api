package routes

import (
	"car-service/cmd/api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupCarRoutes(router *gin.RouterGroup, carController controllers.CarController) {
	cars := router.Group("/cars")
	{
		cars.POST("", carController.CreateCar)
		cars.GET("", carController.GetCars)
	}
}
