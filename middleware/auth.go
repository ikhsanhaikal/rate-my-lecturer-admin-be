package middleware

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	db "github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/mysql"
)

func AuthMiddleware(conn *sql.DB) fiber.Handler {
	return jwtware.New(jwtware.Config{
		JWKSetURLs: []string{os.Getenv("AUTH0_JWKS")},
		SuccessHandler: func(c *fiber.Ctx) error {
			user := c.Locals("token").(*jwt.Token)

			aud, err := user.Claims.GetAudience()

			if err != nil {
				fmt.Printf("%+v\n", err)
				return c.Status(fiber.ErrBadRequest.Code).Send([]byte("meh"))
			}

			fmt.Printf("aud: %+v\n", aud)
			url, err := url.Parse(aud[1])

			if err != nil {
				fmt.Printf("%+v\n", err)
				return c.Status(fiber.ErrBadRequest.Code).Send([]byte("meh"))
			}

			request := http.Request{
				URL: url,
				Header: map[string][]string{
					"Authorization": {fmt.Sprintf("Bearer %s", user.Raw)},
					"Content-type":  {"application/json"},
				},
			}

			request.URL = url
			client := http.Client{}

			resp, err := client.Do(&request)

			if err != nil {
				fmt.Println("failed on fetch")
			}

			userInfo := &UserInfo{}
			body, _ := ioutil.ReadAll(resp.Body)

			if err := json.Unmarshal(body, userInfo); err != nil {
				return c.JSON(fiber.Map{
					"message": "failed on unmarshal",
				})
			}

			queries := db.New(conn)

			// usrRecord, err := queries.GetUserByEmail(context.Background(), userInfo.Email)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					fmt.Println("no user record was found by email")
					queries.CreateUser(c.Context(), db.CreateUserParams{
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

			c.Locals("user", userInfo.Email)

			// fmt.Printf("usrRecord: %+v\n", usrRecord)
			fmt.Printf("body %+v\n", userInfo)

			return c.Next()
		},
		ContextKey: "token",
	})
}
