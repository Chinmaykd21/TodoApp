package utilities

import (
	"net/http"

	"github.com/Chinmaykd21/TodoApp/server/serverErrors"
	"github.com/gofiber/fiber/v2"
)

func ParsingStringToInt(c *fiber.Ctx, id string) (int, error) {
	todoInt, err := c.ParamsInt("id")

	// if there is an error then return
	if err != nil {
		errResponse := serverErrors.New(serverErrors.ParseInt, "")
		c.Status(http.StatusUnprocessableEntity)
		_, err = c.WriteString(errResponse.Error())
		return todoInt, err
	}

	return todoInt, err
}