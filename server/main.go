package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	TodoId      int    `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	IsCompleted bool   `json:"isCompleted"`
}

func main() {
	// load environment variables from env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not read from the env file")
	}

	// an array of Todo
	todos := []Todo{}

	// get the mongoDB connection string
	URI := os.Getenv("MONGO_CONNECTION_URL")
	if URI == "" {
		// TODO: should we use panic & recovery instead of log.fatal?
		log.Fatal("MongoDB connection string cannot be empty")
	}

	// We wait for 5 seconds before we throw a timeout error
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	// connect to mongodb client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		// TODO: should we use panic & recovery instead of log.fatal?
		log.Fatal("Cannot connect to the mongoDB", err)
	}
	fmt.Println("Sucessfully connected to mongoDB atlas")

	// Try a ping to the client
	if err := client.Ping(ctx, nil); err != nil {
		// TODO: should we use panic & recovery instead of log.fatal?
		fmt.Println("Cannot ping the db connection", err)
	}

	// We will disconnect our client instance at the end of the main function.
	defer client.Disconnect(ctx)

	// Create a database in mongoDB atlas
	projectsDatabase := client.Database("projects")

	// Create a collection for todo in database projects
	collectionTodos := projectsDatabase.Collection("todos")

	// Create the fiber app instance
	app := fiber.New()

	// Using middleware to solve CORS issue
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("ALLOWED_ORIGIN_CLIENT"),
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// To return all the posts that are available in our collection
	app.Get("/api/todos", func(c *fiber.Ctx) error {

		cursor, err := collectionTodos.Find(ctx, bson.M{})
		if err != nil {
			return err
		}

		if err = cursor.All(ctx, &todos); err != nil {
			errResponse := serverErrors.New(serverErrors.RetreivalError, "")
			// c.Status(http.StatusUnprocessableEntity)
			_, err = c.WriteString(errResponse.Error())
			return err
		}

		return c.JSON(todos)
	})

	// Add new todo list
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		// create a todo
		todo := &Todo{}

		// try to parse that todo, if it returns an error, return that error
		if err := c.BodyParser(todo); err != nil {
			errResponse := serverErrors.New(serverErrors.BodyParse, "")
			c.Status(http.StatusUnprocessableEntity)
			_, err = c.WriteString(errResponse.Error())
			return err
		}

		// new unique id for the todo
		// idTodo := len(todos) + 1

		// // insert first dummy todo document to the collection
		// _, err := collectionTodos.InsertOne(ctx, bson.D{
		// 	{Key: "id", Value: idTodo},
		// 	{Key: "title", Value: &todo.Title},
		// 	{Key: "body", Value: &todo.Body},
		// 	{Key: "isCompleted", Value: &todo.IsCompleted},
		// })
		// if err != nil {
		// 	// TODO: should we use panic & recovery instead of log.fatal?
		// 	errResponse := serverErrors.New(serverErrors.InsertError, "")
		// 	c.Status(http.StatusUnprocessableEntity)
		// 	_, err = c.WriteString(errResponse.Error())
		// 	return err
		// }
		// fmt.Println("Sucessfully inserted document to the collection todos.")

		// // Query collection to get the latest todos from database
		// // todos = collectionTodos.

		// return new todos
		return c.JSON(todos)
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
