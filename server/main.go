package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Chinmaykd21/TodoApp/server/crudData"
	"github.com/Chinmaykd21/TodoApp/server/customDataStructs"
	"github.com/Chinmaykd21/TodoApp/server/models"
	"github.com/Chinmaykd21/TodoApp/server/utilities"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

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
	// ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	// defer cancel()
	ctx := context.TODO()

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

	// To delete the task which are completed in correctly
	app.Delete("/api/todos/:id?/delete", func(c *fiber.Ctx) error {
		todoId, err := utilities.ParsingStringToInt(c, "id")

		if err != nil {
			return err
		}

		// finding & deleting the record
		errResponse := crudData.DeleteTodo(ctx, c, todoId, collectionTodos)

		if errResponse != nil {
			c.Status(http.StatusNotFound)
			_, err = c.WriteString(errResponse.Error())
			return err
		}

		return c.JSON("Record deleted successfully")
	})

	// To update specific task list from the todos list
	app.Patch("/api/todos/:id?/toggle", func(c *fiber.Ctx) error {

		todoId, err := utilities.ParsingStringToInt(c, "id")

		if err != nil {
			return err
		}

		// Finding & updating the record
		updatedRecord, err := crudData.UpdateTask(ctx, c, todoId, collectionTodos)
		// If there is no error then find this record & update its status
		if err != nil {
			return err
		}

		return c.JSON(updatedRecord)
	})

	// To make server listen on specific port
	PORT := ":" + os.Getenv("SERVER_PORT")
	log.Fatal(app.Listen(PORT))
}
