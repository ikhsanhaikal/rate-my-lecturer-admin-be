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
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()

	conn, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:55001)/playground?parseTime=true")

	if err != nil {
		fmt.Printf("err(sql.Open):%+v\n", err)
		os.Exit(1)
	}

	builder := gql.TypeBuilder{
		DB: conn,
	}

	resolver := gql.Resolver{
		DB: conn,
	}

	labType := builder.LabType()
	lecturerType := builder.LecturerType(labType)
	userType := builder.UserType()

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"lecturers": &graphql.Field{
				Type:    graphql.NewList(lecturerType),
				Resolve: resolver.ListLecturers,
			},
			"lecturer": &graphql.Field{
				Type: lecturerType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: resolver.GetLecturerById,
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"lecturer": &graphql.Field{
				Type:        lecturerType,
				Description: "create a lecturer",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(gql.CreateLecturerInput),
					},
				},
				Resolve: resolver.CreateLecturer,
			},
			"users": &graphql.Field{
				Type:        userType,
				Description: "create a user",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(gql.CreateUserInput),
					},
				},
				Resolve: resolver.CreateUser,
			},
		},
	})

	var AppSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	if err := godotenv.Load(); err != nil {
		os.Exit(1)
	}

	h := handler.New(&handler.Config{
		Schema:     &AppSchema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	app.Use("/graphql", adaptor.HTTPHandler(h))

	app.Listen("127.0.0.1:8080")
}
