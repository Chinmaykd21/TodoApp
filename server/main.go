package main

import (
	"context"
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

	// We will disconnect our client instance at the end of the main function.
	defer client.Disconnect(ctx)

	// TODO: Need to move these API calls to their own packages.
	// To return all the posts that are available in our collection
	app.Get("/api/todos", func(c *fiber.Ctx) error {

		// call function to get the records from the collection
		obtainedTodos, err := crudData.GetAllTodos(ctx, c, todos, collectionTodos)

		if err != nil {
			return err
		}

		return c.JSON(obtainedTodos)
	})

	// TODO: Make this function more efficient. Doing 2 calls to the same collection to get
	// the data seems extremely inefficient

	// Add new todo list
	app.Post("/api/todos", func(c *fiber.Ctx) error {

		todosBefore, errBeforeNewTodo := crudData.GetAllTodos(ctx, c, todos, collectionTodos)
		if errBeforeNewTodo != nil {
			return errBeforeNewTodo
		}

		err := crudData.AddToDo(ctx, c, *todosBefore, collectionTodos)

		if err != nil {
			return err
		}

		todosAfter, errAfterNewTodo := crudData.GetAllTodos(ctx, c, todos, collectionTodos)
		if errAfterNewTodo != nil {
			return errAfterNewTodo
		}

		return c.JSON(todosAfter)
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
