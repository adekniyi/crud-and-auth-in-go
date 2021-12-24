package routes

import (
	controller "goauth/controllers"

	"github.com/gin-gonic/gin"
)

func BookRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("users/books", controller.GetBooks())
	incomingRoutes.POST("users/book", controller.CreateBook())
	incomingRoutes.GET("user/book/:bookId", controller.GetBook())
}
