package utilities

import (
	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/gofiber/fiber/v2"
)

func ParsingStringToInt(c *fiber.Ctx, id string) (int, error) {
	todoInt, err := c.ParamsInt("id")

	// if there is an error then return custom error ParseInt
	if err != nil {
		errResponse := serverErrors.New(serverErrors.ParseInt, err.Error())
		return 0, errResponse
	}

	return todoInt, nil
}
