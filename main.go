package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	routers "github.com/kaar20/todo/Routers"
)

func main() {
	fmt.Println("Welcome to the Todo List !") // Print the message "Hello, World!" to the console
	port := os.Getenv("port")

	if port == "" {
		port = "8080"
	}
	router := gin.New()
	router.Use(gin.Logger())

	routers.TodoRoutes(router)

	router.Run(":" + port)
}
