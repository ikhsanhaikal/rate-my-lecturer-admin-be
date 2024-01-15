package router

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/mydb"
)

type User struct {
	Id    uint   `json:"id,omitempty"`
	Name  string `json:"username"`
	Email string `json:"email"`
}

func UsersRouter(db *sql.DB) *fiber.App {
	router := fiber.New()

	router.Get("/", func(c *fiber.Ctx) error {
		queries := mydb.New(db)
		result, err := queries.ListUsers(context.Background())

		if err != nil {
			return c.JSON(fiber.Map{
				"result": nil,
				"errors": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"result": result,
			"errors": nil,
		})
	})

	router.Get("/:id", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{})
	})

	return router
}
