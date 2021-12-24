package routes

import (
	controller "goauth/controllers"
	"goauth/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoute *gin.Engine) {
	incomingRoute.Use(middleware.Authenticate())
	incomingRoute.GET("/users", controller.GetUsers())
	incomingRoute.GET("/users/:user_id", controller.GetUser())
}
