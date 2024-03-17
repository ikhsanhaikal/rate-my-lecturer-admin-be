package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/gql"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()

	conn, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:55003)/rate_my_lecturer?parseTime=true")

	if err != nil {
		fmt.Printf("err(sql.Open):%+v\n", err)
		os.Exit(1)
	}

	resolver := gql.Resolver{
		DB: conn,
	}

	gqlFactory := gql.GqlFactory{
		DB: conn,
	}

	labType := gqlFactory.LabType()
	lecturerType := gqlFactory.LecturerType(labType)
	userType := gqlFactory.UserType()
	subjectType := gqlFactory.SubjectType()

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"lecturers": &graphql.Field{
				Type: graphql.NewList(lecturerType),
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: resolver.ListLecturers,
			},
			"lecturers_by_pk": &graphql.Field{
				Type: lecturerType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: resolver.GetLecturerById,
			},
			"lecturers_by_lab": &graphql.Field{
				Type: graphql.NewList(lecturerType),
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: resolver.GetLecturersByLab,
			},
			"labs": &graphql.Field{
				Type: graphql.NewList(labType),
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: resolver.ListLabs,
			},
			"labs_by_pk": &graphql.Field{
				Type: labType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: resolver.GetLabById,
			},
			"subjects": &graphql.Field{
				Type: graphql.NewList(subjectType),
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 12,
					},
					"page": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 1,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					queries := sqlcdb.New(conn)
					result, err := queries.ListSubjects(p.Context)
					if err != nil {
						return nil, err
					}
					return result, nil
				},
			},
			"me": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var queries = sqlcdb.New(conn)

					fmt.Printf("ME QUERY CALLED\n")

					var userEmail, ok = p.Context.Value("user").(string)
					if !ok {
						fmt.Printf("err on cast userInfo\n")
						return nil, errors.New("server failed :)")
					}

					var user, err = queries.GetUserByEmail(p.Context, userEmail)

					if err != nil {
						return nil, err
					}

					fmt.Printf("user: %+v\n", user)
					return user, nil
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"create_lecturers_one":   gqlFactory.CreateLecturer(lecturerType),
			"delete_lecturers_by_pk": gqlFactory.DeleteLecturersByPk(lecturerType),
			"delete_lecturers":       gqlFactory.DeleteLecturers(lecturerType),
			"update_lecturers_one":   gqlFactory.UpdateLecturer(lecturerType),
			"create_labs_one":        gqlFactory.CreateLab(labType),
			"delete_labs_by_pk":      gqlFactory.DeleteLabsByPk(labType),
			"user": &graphql.Field{
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

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// app.Use(middleware.AuthMiddleware(conn))

	app.All("/graphql", adaptor.HTTPHandler(h))

	app.Listen("127.0.0.1:8080")
}
