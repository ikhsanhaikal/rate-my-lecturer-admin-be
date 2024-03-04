package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/gql"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/middleware"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()

	conn, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:55001)/playground?parseTime=true")

	if err != nil {
		fmt.Printf("err(sql.Open):%+v\n", err)
		os.Exit(1)
	}

	repository := gql.Repository{
		Db: conn,
	}

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"Lecturers": &graphql.Field{
				Type:    graphql.NewList(repository.LecturerType()),
				Resolve: repository.ListLecturersResolver,
			},
		},
	})

	var AppSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: nil,
	})

	if err := godotenv.Load(); err != nil {
		os.Exit(1)
	}

	app.Use("/api", middleware.AuthMiddleware(conn))

	h := handler.New(&handler.Config{
		Schema:   &AppSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	app.Use("/graphql", adaptor.HTTPHandler(h))

	app.Listen("127.0.0.1:8080")
}
