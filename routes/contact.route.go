package routes

import (
	"example/web-service-gin/controllers"
	"example/web-service-gin/middlewares"

	"github.com/gin-gonic/gin"
)

type ContactRoutes struct {
	contactController controllers.ContactController
}

func NewRouteContact(contactController controllers.ContactController) ContactRoutes {
	return ContactRoutes{contactController}
}

func (cr *ContactRoutes) ContactRoute(rg *gin.RouterGroup) {

	routerGroup := rg.Group("contacts")
	routerGroup.Use(middlewares.CheckAuth(cr.contactController.TokenMaker))
	routerGroup.POST("/", cr.contactController.CreateContact)
	routerGroup.GET("/", cr.contactController.GetAllContacts)
	routerGroup.PATCH("/:contactId", cr.contactController.UpdateContact)
	routerGroup.GET("/:contactId", cr.contactController.GetContactById)
	routerGroup.DELETE("/:contactId", cr.contactController.DeleteContactById)
}
