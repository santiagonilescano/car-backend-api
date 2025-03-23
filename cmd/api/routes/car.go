package routes

import (
	"car-service/cmd/api/controllers"

	"github.com/gin-gonic/gin"
)

// CarRoutes configura las rutas relacionadas con cars
func SetupCarRoutes(router *gin.RouterGroup, carController controllers.CarController) {
	cars := router.Group("/cars")
	{
		cars.POST("", carController.CreateCar)
		cars.GET("", carController.GetCars)
		// Aquí irán más rutas relacionadas con cars:
		// cars.GET("", carHandler.ListCars)
		// cars.GET("/:id", carHandler.GetCar)
		// cars.PUT("/:id", carHandler.UpdateCar)
		// cars.DELETE("/:id", carHandler.DeleteCar)
	}
}
