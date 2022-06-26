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
	idTodo := len(todos) + 1

	// insert first dummy todo document to the collection
	_, err := collectionTodos.InsertOne(ctx, bson.D{
		{Key: "id", Value: idTodo},
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
	fmt.Println("Sucessfully inserted document to the collection todos.")

	return err
}

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
