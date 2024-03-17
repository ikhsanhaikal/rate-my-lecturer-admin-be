package gql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
)

func (factory *GqlFactory) CreateLecturer(returnType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        returnType,
		Description: "create a lecturer",
		Args: graphql.FieldConfigArgument{
			"input": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(CreateLecturersInput),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var queries = sqlcdb.New(factory.DB)

			input, ok := p.Args["input"].(map[string]interface{})

			if !ok {
				fmt.Printf("err can't cast args to struct\n")
				return nil, nil
			}

			fmt.Printf("input: %+v\n", input)

			id, err := queries.CreateLecturer(p.Context, sqlcdb.CreateLecturerParams{
				Name:        input["name"].(string),
				Email:       input["email"].(string),
				Description: utilStringExist(input["description"]),
				Labid:       int32(input["labId"].(int)),
				Gender:      utilStringExist(input["gender"]),
			})

			if err != nil {
				var mysqlErr *mysql.MySQLError
				if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
					return nil, errors.New("data already exists")
				}
				fmt.Printf("failed oh failed\n")
				return nil, err
			}
			data, _ := queries.GetLecturersByPk(p.Context, int32(id))

			return data, nil
		},
	}
}

var UpdateLecturerInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdateLecturerInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"id": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"name": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"email": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"description": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"labId": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
	},
})

func utilStringExist(givenString interface{}) sql.NullString {
	if str, ok := givenString.(string); ok {
		return sql.NullString{
			String: str,
			Valid:  true,
		}
	}
	return sql.NullString{}
}
func utilIntExist(givenInt interface{}) sql.NullInt32 {
	if value, ok := givenInt.(int); ok {
		return sql.NullInt32{
			Int32: int32(value),
			Valid: true,
		}
	}
	return sql.NullInt32{}
}
func (factory *GqlFactory) UpdateLecturer(returnType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type: returnType,
		Args: graphql.FieldConfigArgument{
			"input": &graphql.ArgumentConfig{
				Type: UpdateLecturerInput,
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			queries := sqlcdb.New(factory.DB)

			input, _ := p.Args["input"].(map[string]interface{})
			targetId, _ := input["id"].(int)
			var err = queries.UpdateLecturer(p.Context, sqlcdb.UpdateLecturerParams{
				ID:          int32(targetId),
				Name:        utilStringExist(input["name"]),
				Email:       utilStringExist(input["name"]),
				Description: utilStringExist(input["description"]),
				Gender:      utilStringExist(input["gender"]),
				Labid:       utilIntExist(input["labId"]),
			})

			if err != nil {
				return nil, err
			}

			var lecturer, _ = queries.GetLecturersByPk(p.Context, int32(targetId))

			return lecturer, nil
		},
	}
}

func (factory *GqlFactory) DeleteLecturersByPk(returnType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        returnType,
		Description: "delete a lecturer",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// queries := sqlcdb.New(fac.DB)

			return nil, nil
		},
	}
}

func (factory *GqlFactory) DeleteLecturers(returnType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewList(graphql.Int),
		Args: graphql.FieldConfigArgument{
			"ids": &graphql.ArgumentConfig{
				Type: graphql.NewList(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {

			ids, ok := p.Args["ids"].([]interface{})
			if !ok {
				fmt.Printf("dude it's not okay for some reason to type assertion to []int\n")
				return nil, errors.New("invalid ids")
			}

			fmt.Printf("T of args: %T\n", p.Args["ids"])
			fmt.Printf("ok: %+v\n", ok)

			targets := []string{}

			for _, id := range ids {
				value, ok := id.(int)
				if !ok {
					fmt.Printf("%t\n", ok)
				} else {
					targets = append(targets, fmt.Sprint(value))
				}
			}
			fmt.Printf("targets: %+v\n", targets)
			tx, err := factory.DB.BeginTx(p.Context, nil)

			if err != nil {
				fmt.Printf("err: %+v\n", err)
				return nil, err
			}

			str := strings.Join(targets, ",")

			fmt.Printf("str: %s\n", str)

			_, err = tx.Exec(fmt.Sprintf("DELETE FROM lecturers WHERE lecturers.id IN (%s)", str))

			if err != nil {
				return nil, err
			}
			if err := tx.Commit(); err != nil {
				fmt.Printf("err: on commit %+v\n", err)
				return nil, err
			}

			return targets, nil
		},
	}
}
