package crudData

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Chinmaykd21/TodoApp/server/customDataStructs"
	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllDocuments(ctx context.Context, c *fiber.Ctx, todos []customDataStructs.Todo, collectionTodos *mongo.Collection) (*[]customDataStructs.Todo, error) {
	cursor, err := collectionTodos.Find(ctx, bson.M{})
	if err != nil {
		return &todos, err
	}

	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &todos); err != nil {
		errResponse := serverErrors.New(serverErrors.RetreivalError, "")
		// c.Status(http.StatusUnprocessableEntity)
		_, err = c.WriteString(errResponse.Error())
		return &todos, err
	}

	return &todos, err
}

func AddTodoDocument(ctx context.Context, c *fiber.Ctx, todos []customDataStructs.Todo, collectionTodos *mongo.Collection) error {
	// create an empty todo
	todo := &customDataStructs.Todo{}

	// If it returns an error, return that error
	if err := c.BodyParser(todo); err != nil {
		errResponse := serverErrors.New(serverErrors.BodyParse, "")
		c.Status(http.StatusUnprocessableEntity)
		_, err = c.WriteString(errResponse.Error())
		return err
	}

	// new unique id for the todo
	todoId := len(todos) + 1

	// insert first dummy todo document to the collection
	_, err := collectionTodos.InsertOne(ctx, bson.D{
		{Key: "todoId", Value: todoId},
		{Key: "title", Value: &todo.Title},
		{Key: "body", Value: &todo.Body},
		{Key: "isCompleted", Value: &todo.IsCompleted},
	})

	if err != nil {
		// TODO: should we use panic & recovery instead of log.fatal?
		errResponse := serverErrors.New(serverErrors.InsertError, "")
		c.Status(http.StatusUnprocessableEntity)
		_, err = c.WriteString(errResponse.Error())
		return err
	}

	return err
}

func UpdateTask(ctx context.Context, c *fiber.Ctx, todoId int, collectionTodos *mongo.Collection) (*[]customDataStructs.Todo, error) {
	// Initialize an empty todo
	var todoToUpdate customDataStructs.Todo
	var updatedTodoInSlice []customDataStructs.Todo

	// Define filter condition to find the document with proper val
	filter := bson.D{{Key: "todoId", Value: bson.D{{Key: "$eq", Value: todoId}}}}

	// Find and store the id whose status is to be toggled to struct todo
	err := collectionTodos.FindOne(ctx, filter).Decode(&todoToUpdate)
	if err != nil {
		return &updatedTodoInSlice, err
	}

	// Toggle its current state & store it
	toggleIsCompleted := !todoToUpdate.IsCompleted

	// Update condition
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "isCompleted", Value: toggleIsCompleted}}}}

	// Update the same ID with the new state
	updateResult, err := collectionTodos.UpdateOne(ctx, filter, update)

	if err != nil {
		fmt.Println("There was an error while updating the document")
		return &updatedTodoInSlice, err
	}

	if updateResult.MatchedCount != updateResult.ModifiedCount {
		fmt.Println("Modified & Matched record count do not match")
		return &updatedTodoInSlice, err
	}

	updatedTodoInSlice = append(updatedTodoInSlice, todoToUpdate)

	return &updatedTodoInSlice, err
}
