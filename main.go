package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	_ "github.com/go-sql-driver/mysql"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/mydb"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/router"
	"github.com/joho/godotenv"
)

func AuthMiddleware(db *sql.DB) fiber.Handler {
	return jwtware.New(jwtware.Config{
		JWKSetURLs: []string{os.Getenv("AUTH0_JWKS")},
		SuccessHandler: func(c *fiber.Ctx) error {
			user := c.Locals("user").(*jwt.Token)

			aud, err := user.Claims.GetAudience()

			if err != nil {

			}

			fmt.Printf("aud: %+v\n", aud)
			url, err := url.Parse(aud[1])

			request := http.Request{
				URL: url,
				Header: map[string][]string{
					"Authorization": {fmt.Sprintf("Bearer %s", user.Raw)},
					"Content-type":  {"application/json"},
				},
			}

			type UserInfo struct {
				Name    string `json:"name"`
				Picture string `json:"picture"`
				Email   string `json:"email"`
			}

			request.URL = url
			client := http.Client{}

			resp, err := client.Do(&request)

			if err != nil {
				fmt.Println("failed on fetch")
			}

			userInfo := &UserInfo{}
			body, err := ioutil.ReadAll(resp.Body)

			if err := json.Unmarshal(body, userInfo); err != nil {
				return c.JSON(fiber.Map{
					"message": "failed on unmarshal",
				})
			}

			queries := mydb.New(db)

			usrRecord, err := queries.GetUserByEmail(context.Background(), userInfo.Email)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					fmt.Println("no user record was found by email")
					queries.CreateUser(c.Context(), mydb.CreateUserParams{
						Email: userInfo.Email,
						Name:  userInfo.Name,
					})
				} else {
					return c.JSON(fiber.Map{
						"error":   err,
						"message": "failed for some reaseon",
					})
				}
			}

			fmt.Printf("usrRecord: %+v\n", usrRecord)
			fmt.Printf("body %+v\n", userInfo)

			return c.Next()
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.JSON(fiber.Map{
				"error":   err.Error(),
				"message": "no token / invalid token",
			})
		},
	})
}

func main() {
	app := fiber.New()

	if err := godotenv.Load(); err != nil {
		os.Exit(1)
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:55001)/playground?parseTime=true")

	if err != nil {
		fmt.Printf("err(sql.Open):%+v\n", err)
		os.Exit(1)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "ok",
			"status":  200,
		})
	}).Name("index")

	app.Use("/todos", AuthMiddleware(db))

	app.Mount("/users", router.UsersRouter(db))

	app.Listen("127.0.0.1:3000")
}
