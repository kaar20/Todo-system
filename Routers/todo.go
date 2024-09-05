package routers

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/kaar20/todo/Controllers"
)

func TodoRoutes(router *gin.Engine) {
	router.GET("todo/list", controllers.ListTodo())
	router.GET("todo/:id", controllers.GetTodo())
	router.POST("todo/add", controllers.AddTodo())
	router.PATCH("todo/:id", controllers.UpdateTodo())
	router.DELETE("todo/:id", controllers.DeleteTodo())
}
