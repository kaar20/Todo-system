package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kaar20/todo/database"
	"github.com/kaar20/todo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// router.GET("todo/list", controller.ListTodo())
// router.POST("todo/add", controller.AddTodo())
// router.PUT("todo/:id", controller.UpdateTodo())
// router.DELETE("todo/:id", controller.DeleteTodo())
var todoCollection *mongo.Collection = database.OpenCollection(database.Client, "todoList")
var validate = validator.New()

// Get all Todo List
func ListTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var todoList []bson.M

		findResult, findError := todoCollection.Find(ctx, bson.M{})
		defer cancel()
		if findError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": findError.Error()})
			return
		}

		err := findResult.All(ctx, &todoList)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, todoList)

		// var todoList []bson.M
		// findResult,Err := todoCollection.Find(,bson.M{}).Decode(&todoList)

	}
}

func GetTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		todoId := c.Param("id")
		var foundTodo models.TodoModel

		err := todoCollection.FindOne(ctx, bson.M{"todo_id": todoId}).Decode(&foundTodo)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}
		c.JSON(http.StatusOK, foundTodo)

	}
}

// add to to the list of todos

func AddTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		// defer cancel()
		var todo models.TodoModel
		if err := c.BindJSON(&todo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": err.Error(),
			})
		}
		validation := validate.Struct(todo)
		if validation != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": validation.Error()})
		}
		todo.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		todo.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		todo.ID = primitive.NewObjectID()
		todo.Todo_id = todo.ID.Hex()
		todo.Is_completed = false

		insertResult, InsertError := todoCollection.InsertOne(ctx, todo)
		defer cancel()
		if InsertError != nil {
			c.JSON(http.StatusMisdirectedRequest, gin.H{"InsertError": InsertError.Error()})
		}
		c.JSON(http.StatusAccepted, gin.H{"Inserted Data": insertResult})

	}
}

// Update To the todo
func UpdateTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel() // Ensure the context is always canceled

		// Get the todo ID from the request parameters
		todoId := c.Param("id")

		// Find the existing todo
		var todoFound models.TodoModel
		err := todoCollection.FindOne(ctx, bson.M{"todo_id": todoId}).Decode(&todoFound)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		// Bind and validate the updated todo
		var updatedTodo models.TodoModel
		if err := c.BindJSON(&updatedTodo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}

		// Validate the updated fields
		if validation := validate.Struct(updatedTodo); validation != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": validation.Error()})
			return
		}
		updatedTodo.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		// Build the update document dynamically
		update := bson.M{"$set": bson.M{"Updated_at": updatedTodo.Updated_at}}
		if updatedTodo.Title != "" {
			update["$set"].(bson.M)["Title"] = updatedTodo.Title
		}
		if updatedTodo.Description != "" {
			update["$set"].(bson.M)["Description"] = updatedTodo.Description
		}
		if updatedTodo.Is_completed {
			update["$set"].(bson.M)["Is_completed"] = updatedTodo.Is_completed
		}
		if updatedTodo.User != "" {
			update["$set"].(bson.M)["User"] = updatedTodo.User
		}

		// Perform the update operation
		filter := bson.M{"todo_id": todoId}
		_, err = todoCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Updated Todo": todoFound})
	}
}

//Delete Todo

func DeleteTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var todoId = c.Param("id")
		filter := bson.M{"todo_id": todoId}

		result, err := todoCollection.DeleteOne(ctx, filter)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Deleted Todo ID": todoId})

	}
}
