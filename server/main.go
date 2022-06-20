package main

import (
	"log"
	"net/http"

	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	IsCompleted bool   `json:"isCompleted"`
}

func main() {
	// an array of Todo
	todos := []Todo{}

	app := fiber.New()

	// using middleware to solve CORS issue
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// To return all the posts that are available
	app.Get("/api/todos", func(c *fiber.Ctx) error {
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

		// otherwise, give a unique id to that todo
		todo.Id = len(todos) + 1

		// append that todo to todos array
		todos = append(todos, *todo)

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
			if todo.Id == id {
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
	log.Fatal(app.Listen(":4000"))
}
