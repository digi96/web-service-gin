package routes

import (
	"example/web-service-gin/controllers"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	authController controllers.AuthController
}

func NewRouteAuth(authController controllers.AuthController) AuthRoutes {
	return AuthRoutes{authController}
}

func (ar *AuthRoutes) AuthRoute(rg *gin.RouterGroup) {
	routerGroup := rg.Group("auth")
	routerGroup.GET("/login/:contactId", ar.authController.Login)

}
