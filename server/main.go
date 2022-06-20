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

func main() {
	// an array of Todo
	todos := []Todo{}

	app := fiber.New()

	// using middleware to solve CORS issue
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// To check the health of the server
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Home route
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.JSON(todos)
	// })

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
			return err
		}

		// otherwise, give a unique id to that todo
		todo.Id = len(todos) + 1

		// append that todo to todos array
		todos = append(todos, *todo)

		// return new todos
		return c.JSON(todos)
	})

	// To get specific id from the url
	app.Patch("api/todos/:id?/toggle", func(c *fiber.Ctx) error {
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
				break
			}
		}

		return c.JSON(todos)
	})

	// To make server listen on specific port
	log.Fatal(app.Listen(":4000"))
}
