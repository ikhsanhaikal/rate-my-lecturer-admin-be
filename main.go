package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/handler"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/gql"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
	"github.com/joho/godotenv"
	"gopkg.in/guregu/null.v3"
)

type InputCreateCourse struct {
	Lecturer int
	Subject  int
	Year     time.Time
	Semester string
}

func main() {
	app := fiber.New()
	if err := godotenv.Load(); err != nil {
		os.Exit(1)
	}

	dbUrl := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("root:password@tcp(%s)/%s?parseTime=true", dbUrl, dbName)
	fmt.Println(connStr)
	conn, err := sql.Open("mysql", connStr)

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
	traitType := gqlFactory.TraitType()
	lecturerType := gqlFactory.LecturerType(labType, traitType)
	userType := gqlFactory.UserType()
	subjectType := gqlFactory.SubjectType()
	courseType := gqlFactory.CourseType(subjectType, lecturerType)
	reviewType := gqlFactory.ReviewType(courseType, traitType, userType)
	foo := gqlFactory.ListLecturersResult(lecturerType)
	sortType := gqlFactory.SortType()

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"lecturers": &graphql.Field{
				Type: foo,
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"sort": &graphql.ArgumentConfig{
						Type: sortType,
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
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "LecturersByLabResult",
					Fields: graphql.Fields{
						"data": &graphql.Field{
							Type: graphql.NewList(lecturerType),
						},
						"total": &graphql.Field{
							Type: graphql.Int,
						},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"limit": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: resolver.GetLecturersByLab,
			},
			"labs": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "ListLabsResult",
					Fields: graphql.Fields{
						"data": &graphql.Field{
							Type: graphql.NewList(labType),
						},
						"total": &graphql.Field{
							Type: graphql.Int,
						},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"sort": &graphql.ArgumentConfig{
						Type: sortType,
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
			"reviews_by_lecturer": &graphql.Field{
				Name: "reviews_by_lecturer",
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "ReviewsByLecturerResult",
					Fields: graphql.Fields{
						"data": &graphql.Field{
							Type: graphql.NewList(reviewType),
						},
						"total": &graphql.Field{
							Type: graphql.Int,
						},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Description: "id of the lecturer",
						Type:        graphql.NewNonNull(graphql.Int),
					},
					"limit": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: gqlFactory.GetReviewsByLecturer,
			},
			"courses_by_lecturer": &graphql.Field{
				Name: "courses_by_lecturer",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"limit": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "ListCoursesByLecturerResult",
					Fields: graphql.Fields{
						"data": &graphql.Field{
							Type: graphql.NewList(courseType),
						},
						"total": &graphql.Field{
							Type: graphql.Int,
						},
					},
				}),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					lecturerId, ok := p.Args["id"].(int)
					limit, _ := p.Args["limit"].(int)
					page, _ := p.Args["page"].(int)
					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}
					queries := sqlcdb.New(gqlFactory.DB)
					result, err := queries.ListCourseByLecturer(p.Context, sqlcdb.ListCourseByLecturerParams{
						Lecturerid: sql.NullInt64{Int64: int64(lecturerId), Valid: true},
						Limit:      int32(limit),
						Offset:     int32(limit) * int32(page-1),
					})

					if err != nil {
						return nil, err
					}

					total, err := queries.CountCourseByLecturer(p.Context, sql.NullInt64{Int64: int64(lecturerId), Valid: true})

					if err != nil {
						return nil, err
					}

					return struct {
						Data  []sqlcdb.Class
						Total int64
					}{
						Data:  result,
						Total: total,
					}, nil
				},
			},
			"courses_by_pk": &graphql.Field{
				Type: courseType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					courseId, ok := p.Args["id"].(int)
					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					queries := sqlcdb.New(gqlFactory.DB)

					course, err := queries.GetCourseById(p.Context, int64(courseId))

					if err != nil {
						return nil, err
					}

					return course, nil
				},
			},
			"subjects": &graphql.Field{
				Type: graphql.NewObject(graphql.ObjectConfig{
					Name: "ListSubjectsResult",
					Fields: graphql.Fields{
						"data": &graphql.Field{
							Type: graphql.NewList(subjectType),
						},
						"total": &graphql.Field{
							Type: graphql.Int,
						},
					},
				}),
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"sort": &graphql.ArgumentConfig{
						Type: sortType,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					limit, _ := p.Args["limit"].(int)
					page, _ := p.Args["page"].(int)

					queries := sqlcdb.New(conn)

					result, err := queries.ListSubjects(p.Context, sqlcdb.ListSubjectsParams{
						Limit:  int32(limit),
						Offset: int32(limit) * int32(page-1),
					})

					if err != nil {
						return nil, err
					}

					total, err := queries.CountSubjects(p.Context)

					if err != nil {
						return nil, err
					}

					return struct {
						Data  []sqlcdb.Subject
						Total int64
					}{
						Data:  result,
						Total: total,
					}, nil
				},
			},
			"subjects_by_pk": &graphql.Field{
				Type: subjectType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.Int),
						Description: "subject id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					subjectId, ok := p.Args["id"].(int)

					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					queries := sqlcdb.New(gqlFactory.DB)

					subject, err := queries.GetSubjectsByPk(p.Context, int64(subjectId))

					if err != nil {
						return nil, err
					}

					return subject, nil
				},
			},
			"me": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					var queries = sqlcdb.New(conn)

					fmt.Printf("ME QUERY CALLED\n")

					var userEmail, ok = p.Context.Value("editor").(string)
					if !ok {
						fmt.Printf("err on cast userInfo\n")
						return nil, errors.New("server failed :)")
					}

					var user, err = queries.GetEditorByEmail(p.Context, userEmail)

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
			"create_lecturers":       gqlFactory.CreateLecturer(lecturerType),
			"delete_lecturers_by_pk": gqlFactory.DeleteLecturersByPk(lecturerType),
			"delete_lecturers":       gqlFactory.DeleteLecturers(lecturerType),
			"update_lecturers_by_pk": gqlFactory.UpdateLecturer(lecturerType),
			"update_courses_by_pk": &graphql.Field{
				Type: courseType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
							Name: "UpdateCoursesInput",
							Fields: graphql.InputObjectConfigFieldMap{
								"id": &graphql.InputObjectFieldConfig{
									Type:        graphql.NewNonNull(graphql.Int),
									Description: "course id",
								},
								"subject": &graphql.InputObjectFieldConfig{
									Type: graphql.Int,
								},
								"semester": &graphql.InputObjectFieldConfig{
									Type: graphql.String,
								},
								"year": &graphql.InputObjectFieldConfig{
									Type: graphql.DateTime,
								},
							},
						})),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap, ok := p.Args["input"].(map[string]interface{})

					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					b, err := json.Marshal(inputMap)

					if err != nil {
						return nil, err
					}

					updateCourseInput := &struct {
						Id       *int
						Subject  *int
						Semester *string
						Year     *time.Time
					}{}

					if err := json.Unmarshal(b, updateCourseInput); err != nil {
						return nil, err
					}

					queries := sqlcdb.New(gqlFactory.DB)

					subject := sql.NullInt64{}
					year := sql.NullTime{}
					if updateCourseInput.Subject != nil {
						subject.Int64 = int64(*updateCourseInput.Subject)
						subject.Valid = true
					}

					if updateCourseInput.Year != nil {
						year.Time = *updateCourseInput.Year
						year.Valid = true
					}

					if _, err := queries.UpdateCourse(p.Context, sqlcdb.UpdateCourseParams{
						ID:        int64(*updateCourseInput.Id),
						Subjectid: subject,
						Semester:  null.StringFromPtr(updateCourseInput.Semester),
						Year:      year,
					}); err != nil {
						return nil, err
					}

					course, err := queries.GetCourseById(p.Context, int64(*updateCourseInput.Id))

					if err != nil {
						return nil, err
					}

					return course, nil
				},
			},
			"create_labs":       gqlFactory.CreateLab(labType),
			"delete_labs_by_pk": gqlFactory.DeleteLabsByPk(labType),
			"create_courses": &graphql.Field{
				Type:        courseType,
				Description: "create a course and assing / connect to a lecturer by id",
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
							Name: "CreateCoursesInput",
							Fields: graphql.InputObjectConfigFieldMap{
								"year": &graphql.InputObjectFieldConfig{
									Type: graphql.DateTime,
								},
								"semester": &graphql.InputObjectFieldConfig{
									Type: graphql.String,
								},
								"subject": &graphql.InputObjectFieldConfig{
									Type:        graphql.NewNonNull(graphql.Int),
									Description: "subject id",
								},
								"lecturer": &graphql.InputObjectFieldConfig{
									Type:        graphql.NewNonNull(graphql.Int),
									Description: "lecturer id",
								},
							},
						})),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap, ok := p.Args["input"].(map[string]interface{})

					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					b, err := json.Marshal(inputMap)

					if err != nil {
						return nil, err
					}

					input := &InputCreateCourse{}

					if err := json.Unmarshal(b, input); err != nil {
						return nil, err
					}

					queries := sqlcdb.New(gqlFactory.DB)

					id, err := queries.CreateCourse(p.Context, sqlcdb.CreateCourseParams{
						Lecturerid: sql.NullInt64{Int64: int64(input.Lecturer), Valid: true},
						Subjectid:  sql.NullInt64{Int64: int64(input.Subject), Valid: true},
						Semester:   null.NewString(input.Semester, true),
						Year:       sql.NullTime{Time: input.Year, Valid: true},
					})

					if err != nil {
						fmt.Printf("sql error %+v\n", err)
						return nil, err
					}

					c, err := queries.GetCourseById(p.Context, id)

					fmt.Printf("course: %+v\n", c)

					if err != nil {
						return nil, err
					}

					return c, nil
				},
			},
			"create_subjects": &graphql.Field{
				Name: "create_subjects",
				Type: subjectType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewInputObject(graphql.InputObjectConfig{
							Name: "CreateSubjectsInput",
							Fields: graphql.InputObjectConfigFieldMap{
								"name": &graphql.InputObjectFieldConfig{
									Type: graphql.NewNonNull(graphql.String),
								},
								"description": &graphql.InputObjectFieldConfig{
									Type: graphql.String,
								},
							},
						}),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					createSubjectInputMap, ok := p.Args["input"].(map[string]interface{})

					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					b, err := json.Marshal(createSubjectInputMap)

					if err != nil {
						return nil, err
					}

					editorEmail, ok := p.Context.Value("editor").(string)

					fmt.Printf("editorEmail: %s\n", editorEmail)

					if !ok {
						fmt.Printf("err on cast editorInfo\n")
						return nil, errors.New("server failed :)")
					}

					queries := sqlcdb.New(gqlFactory.DB)
					editor, err := queries.GetEditorByEmail(p.Context, editorEmail)

					if err != nil {
						fmt.Printf("CreateSubject Mutation -> GetEditorByEmail err: %+v\n", err)
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					createSubjectInput := &struct {
						Name       string
						Decription *string
					}{}

					if err := json.Unmarshal(b, createSubjectInput); err != nil {
						return nil, err
					}

					subjectId, err := queries.CreateSubject(p.Context, sqlcdb.CreateSubjectParams{
						Name:        createSubjectInput.Name,
						Description: null.StringFromPtr(createSubjectInput.Decription),
						Editorid:    sql.NullInt64{Int64: editor.ID, Valid: true},
					})

					if err != nil {
						return nil, err
					}

					subject, err := queries.GetSubjectsByPk(p.Context, subjectId)

					if err != nil {
						return nil, err
					}

					return subject, nil
				},
			},
			"update_subjects_by_pk": &graphql.Field{
				Type: subjectType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.NewInputObject(graphql.InputObjectConfig{
							Name: "UpdateSubjectsInput",
							Fields: graphql.InputObjectConfigFieldMap{
								"id": &graphql.InputObjectFieldConfig{
									Type:        graphql.NewNonNull(graphql.Int),
									Description: "subject id",
								},
								"name": &graphql.InputObjectFieldConfig{
									Type: graphql.String,
								},
								"description": &graphql.InputObjectFieldConfig{
									Type: graphql.String,
								},
							},
						})),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					inputMap, ok := p.Args["input"].(map[string]interface{})
					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					bites, err := json.Marshal(inputMap)

					if err != nil {
						return nil, err
					}

					updateSubjectInput := &struct {
						SubjectId   int `json:"id"`
						Name        *string
						Description *string
					}{}

					if err := json.Unmarshal(bites, updateSubjectInput); err != nil {
						return nil, err
					}

					queries := sqlcdb.New(gqlFactory.DB)

					if err := queries.UpdateSubject(p.Context, sqlcdb.UpdateSubjectParams{
						ID:          int64(updateSubjectInput.SubjectId),
						Name:        null.StringFromPtr(updateSubjectInput.Name),
						Description: null.StringFromPtr(updateSubjectInput.Description),
					}); err != nil {
						return nil, err
					}

					subject, err := queries.GetSubjectsByPk(p.Context, int64(updateSubjectInput.SubjectId))

					if err != nil {
						return nil, err
					}

					return subject, nil
				},
			},
			"delete_courses_by_pk": &graphql.Field{
				Type: courseType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					courseId, ok := p.Args["id"].(int)

					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					queries := sqlcdb.New(gqlFactory.DB)

					course, err := queries.GetCourseById(p.Context, int64(courseId))

					if err != nil {
						return nil, err
					}

					if _, err := queries.DeleteCourseById(p.Context, int64(courseId)); err != nil {
						return nil, err
					}

					return course, nil
				},
			},
			"delete_subjects_by_pk": &graphql.Field{
				Type: subjectType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.Int),
						Description: "subject id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					subjectId, ok := p.Args["id"].(int)

					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					queries := sqlcdb.New(gqlFactory.DB)
					subject, err := queries.GetSubjectsByPk(p.Context, int64(subjectId))

					if err != nil {
						return nil, err
					}

					if _, err := queries.DeleteSubjectByPk(p.Context, int64(subjectId)); err != nil {
						return nil, err
					}

					return subject, nil
				},
			},
			"delete_reviews_by_pk": &graphql.Field{
				Type: reviewType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.NewNonNull(graphql.Int),
						Description: "review id",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					fmt.Printf("----Delete reviews resolve was called-----\n")
					reviewId, ok := p.Args["id"].(int)

					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					queries := sqlcdb.New(gqlFactory.DB)

					review, err := queries.GetReviewById(p.Context, int64(reviewId))

					if err != nil {
						return nil, err
					}

					if _, err := queries.DeleteReviewById(p.Context, int64(reviewId)); err != nil {
						return nil, err
					}

					return review, nil
				},
			},
		},
	})

	var AppSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})

	h := handler.New(&handler.Config{
		Schema:     &AppSchema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	app.Use(func(c *fiber.Ctx) error {
		// Check if the path starts with /graphql
		if strings.Contains(c.Path(), "graphql") {
			// If the path starts with /graphql, call next middleware
			return c.Next()
		}
		// Serve the index.html file for any other path
		return c.SendFile("./path/to/your/react/app/build/index.html")
	})
	app.Use("/graphql", cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	// app.Use(func(c *fiber.Ctx) error {
	// headers := c.GetReqHeaders()
	// fmt.Printf("request headers: %+v\n", headers)
	// return c.Next()
	// })
	// app.Use(middleware.AuthMiddleware(conn))
	//app.Use(middleware.AuthMiddleware(conn))
	app.All("/graphql", adaptor.HTTPHandler(h))
	//app.Listen("127.0.0.1:8080")
	app.Listen(":8080")
}
