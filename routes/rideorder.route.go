package routes

import (
	"example/web-service-gin/controllers"

	"github.com/gin-gonic/gin"
)

type RideOrderRoutes struct {
	rideOrderController controllers.RideOrderController
}

func NewRouteRideOrder(rideOrderController controllers.RideOrderController) RideOrderRoutes {
	return RideOrderRoutes{rideOrderController}
}

func (ror *RideOrderRoutes) RideOrderRoutes(rg *gin.RouterGroup) {

	router := rg.Group("rideorders")
	router.POST("/", ror.rideOrderController.CreatRideOrder)
	router.GET("/:orderId", ror.rideOrderController.GetOrderById)
	router.PATCH("/:orderId", ror.rideOrderController.UpdateRideOrder)
}
