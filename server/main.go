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
	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// get the mongoDB connection string
	URI := os.Getenv("MONGO_CONNECTION_URL")
	if URI == "" {
		// TODO: should we use panic & recovery instead of log.fatal?
		log.Fatal("MongoDB connection string cannot be empty")
	}

	// We wait for 5 seconds before we throw a timeout error
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	fmt.Println("This is the ctx", ctx)

	// connect to mongodb client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		// TODO: should we use panic & recovery instead of log.fatal?
		log.Fatal("Cannot connect to the mongoDB", err)
	}
	fmt.Println("Sucessfully connected to mongoDB atlas")

	// Try a ping to the client
	// if err := client.Ping(ctx, nil); err != nil {
	// 	// TODO: should we use panic & recovery instead of log.fatal?
	// 	log.Fatal("Cannot ping the db connection", err)
	// }

	// We will disconnect our client instance at the end of the main function.
	defer client.Disconnect(ctx)

	// Create a database in mongoDB atlas
	projectsDatabase := client.Database("projects")

	// Create a collection for todo in database projects
	collectionTodos := projectsDatabase.Collection("todos")

	fmt.Println("Successfully created todos collection in projects database")

	// To return all the posts that are available in our collection
	app.Get("/api/todos", func(c *fiber.Ctx) error {

		// call function to get the records from the collection
		obtainedTodos, err := crudData.GetAllDocuments(ctx, c, todos, collectionTodos)

		if err != nil {
			return err
		}

		return c.JSON(obtainedTodos)
	})

	// Add new todo list
	app.Post("/api/todos", func(c *fiber.Ctx) error {

		err := crudData.AddTodoDocument(ctx, c, todos, collectionTodos)

		if err != nil {
			return err
		}

		obtainedTodos, err := crudData.GetAllDocuments(ctx, c, todos, collectionTodos)
		if err != nil {
			return err
		}

		return c.JSON(obtainedTodos)
	})

	// To delete the task which are completed in correctly
	// app.Delete("/api/todos/:id?/delete")

	// To update specific task list from the todos list
	app.Patch("/api/todos/:id?/toggle", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		// if there is an error then return
		if err != nil {
			errResponse := serverErrors.New(serverErrors.ParseInt, "")
			c.Status(http.StatusUnprocessableEntity)
			_, err = c.WriteString(errResponse.Error())
			return err
		}

		// otherwise iterate through all the todos & update
		// todo when the id is matched
		for index, todo := range todos {
			if todo.TodoId == id {
				todos[index].IsCompleted = !todos[index].IsCompleted
				return c.JSON(todos)
			}
		}

		errResponse := serverErrors.New(serverErrors.InvalidID, "")
		c.Status(http.StatusUnprocessableEntity)
		_, err = c.WriteString(errResponse.Error())
		return err
	})

	// To make server listen on specific port
	PORT := ":" + os.Getenv("SERVER_PORT")
	log.Fatal(app.Listen(PORT))
}
