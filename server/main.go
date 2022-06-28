package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Chinmaykd21/TodoApp/server/crudData"
	"github.com/Chinmaykd21/TodoApp/server/customDataStructs"
	"github.com/Chinmaykd21/TodoApp/server/models"
	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/Chinmaykd21/TodoApp/server/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

// Default time to set for a context to be expired.
const defaultTimeout = 5 * time.Second

func main() {
	// Create the fiber app instance
	app := fiber.New()

	// Using middleware to solve CORS issue
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOWED_ORIGIN_CLIENT"),
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// load environment variables from env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not read from the env file")
	}

	// an empty array of Todo
	todos := []customDataStructs.Todo{}

	// TODO: Why does this line causes "context deadline exceeded" error?
	// The reason being, this context was shared across all the API requests.
	// Each user needs to have their own context with its own timeout, & if
	// that API call does not get a response before the context dies,
	// that request should fail.
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Get the client & todo collection
	client, collectionTodos := models.CreatConnection(ctx)

	// We will disconnect our client instance at the end of the main
	// function.
	defer func() {
		fmt.Println("Closing the mongoDB client connection")
		client.Disconnect(ctx)
	}()

	// TODO: Need to move these API calls to their own packages.

	// API route to get all todos from todos mongoDB collection
	app.Get("/api/todos", func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		// Get all todos from todos mongoDB collection
		obtainedTodos, errResponse := crudData.GetAllTodos(ctx, c, todos, collectionTodos)

		// If there is any error in finding the todos return 404
		// to client
		if errResponse != nil {
			c.Status(http.StatusNotFound)
			_, err = c.WriteString(errResponse.Error())
			return err
		}

		// Return found todos to the client
		return c.JSON(obtainedTodos)
	})

	// TODO: Make this function more efficient. Doing 2 calls to the
	// same collection to get the data seems extremely inefficient

	// API route to add a new todo list to todos mongoDB collection
	// & return all todos to client
	app.Post("/api/todos", func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		// Getting todos before adding new todo to create a unique
		// id for new todo
		todosBefore, errResponse := crudData.GetAllTodos(ctx, c, todos, collectionTodos)

		// If there are any error while retreiving todos, return 404
		// to client
		if errResponse != nil {
			c.Status(http.StatusNotFound)
			_, err = c.WriteString(errResponse.Error())
			return err
		}

		// Add new todo to todos mongoDB collection
		errResponse = crudData.AddTodo(ctx, c, *todosBefore, collectionTodos)

		// If there are any error while inserting new todo to todos
		// mongoDB collection then return a custom error called
		// InsertError & return that to client
		if errResponse != nil {
			c.Status(http.StatusUnprocessableEntity)
			_, err = c.WriteString(errResponse.Error())
			return err
		}

		// Return new todos to client
		return c.JSON(todos)
	})

	// API route to delete matching todo from todos mongoDB collection
	app.Delete("/api/todos/:id?/delete", func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		// Converting string todoId to int
		todoId, errResponse := utilities.ParsingStringToInt(c, "id")

		// If there are any error while doing type conversion return
		// custom error called InvalidID to the client
		if errResponse != nil {
			c.Status(http.StatusUnprocessableEntity)
			_, err = c.WriteString(errResponse.Error())
			return err
		}

		// Use the converted todoId to delete the todo from todos
		// mongoDB collection
		errResponse = crudData.DeleteTodo(ctx, c, todoId, collectionTodos)

		// If there are any error return while deleting the record
		// from mongoDB collection, return 404 error to the client
		if errResponse != nil {
			c.Status(http.StatusNotFound)
			_, err = c.WriteString(errResponse.Error())
			return err
		}

		// return updated todos to client
		return c.JSON(todos)
	})

	// API route to edit specific task list from todos list
	app.Patch("/api/todos/:id?/edit", func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		// Converting string todoId to int
		todoId, errResponse := utilities.ParsingStringToInt(c, "id")

		// If there are any error while doing type conversion return
		// custom error called InvalidID to the client
		if errResponse != nil {
			c.Status(http.StatusUnprocessableEntity)
			_, err = c.WriteString(errResponse.Error())
			return err
		}

		// initializing an empty todo & attaching the todo coming from
		// request body to that todo
		todo := &customDataStructs.Todo{}

		// If it returns an error, return that error
		if err := c.BodyParser(todo); err != nil {
			errResponse := serverErrors.New(serverErrors.BodyParse, "")
			return errResponse
		}

		// Update the todo associate with the todo
		editedTodo, errResponse := crudData.EditTodo(ctx, c, todoId, *todo, todos, collectionTodos)

		// If there is no error then find this record & update its status
		if errResponse != nil {
			c.Status(http.StatusUnprocessableEntity)
			_, err = c.WriteString(errResponse.Error())
			return err
		}

		return c.JSON(editedTodo)
	})

	// API route to update specific task list completed state from the todos list
	app.Patch("/api/todos/:id?/toggle", func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		// Converting string todoId to int
		todoId, errResponse := utilities.ParsingStringToInt(c, "id")

		// If there are any error while doing type conversion return
		// custom error called InvalidID to the client
		if errResponse != nil {
			c.Status(http.StatusUnprocessableEntity)
			_, err = c.WriteString(errResponse.Error())
			return err
		}

		// Use the converted todoId to update the matching todo in
		// todos mongoDB collection
		updatedRecord, errResponse := crudData.UpdateTodo(ctx, c, todoId, collectionTodos)

		// If there is no error then find this record & update its status
		if errResponse != nil {
			c.Status(http.StatusUnprocessableEntity)
			_, err = c.WriteString(errResponse.Error())
			return err
		}

		return c.JSON(updatedRecord)
	})

	// To make server listen on specific port
	PORT := ":" + os.Getenv("SERVER_PORT")
	log.Fatal(app.Listen(PORT))
}
