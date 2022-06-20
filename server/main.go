package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	IsCompleted bool   `json:"isCompleted"`
}

type errMessage struct {
	Code    int    `json:"errCode"`
	Message string `json:"errMessage"`
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
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(todos)
	})

	// Add new todo list
	app.Post("/", func(c *fiber.Ctx) error {
		// create a todo
		todo := &Todo{}

		// try to parse that todo, if it returns an error, return that error
		if err := c.BodyParser(todo); err != nil {
			errResponse := &errMessage{
				Code:    1,
				Message: "Invalid Todo",
			}

			return c.JSON(*errResponse)
		}

		// otherwise, give a unique id to that todo
		todo.Id = len(todos) + 1

		// append that todo to todos array
		todos = append(todos, *todo)

		// return new todos
		return c.JSON(todos)
	})

	// To update specific task list from the todos list
	app.Patch("/:id?/toggle", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		// if there is an error then return
		if err != nil {
			return c.SendString("Invalid id")
		}

		// otherwise iterate through all the todos & update
		// todo when the id is matched
		for index, todo := range todos {
			if todo.Id == id {
				todos[index].IsCompleted = !todos[index].IsCompleted
				return c.JSON(todos)
			}
		}

		errResponse := &errMessage{
			Code:    2,
			Message: "Invalid ID",
		}

		return c.JSON(*errResponse)
	})

	// To make server listen on specific port
	log.Fatal(app.Listen(":4000"))
}
