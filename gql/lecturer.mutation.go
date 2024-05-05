package gql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
	"gopkg.in/guregu/null.v3"
)

type CreateLecturerDto struct {
	Name        string
	Email       string
	Description *string
	Gender      *string
	Labid       int
}

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
			var editorEmail = p.Context.Value("editor").(string)

			inputMap, ok := p.Args["input"].(map[string]interface{})

			if !ok {
				fmt.Printf("err can't cast args to struct\n")
				return nil, nil
			}

			b, err := json.Marshal(inputMap)

			if err != nil {
				return nil, err
			}

			createLecturerDto := &CreateLecturerDto{}

			if err := json.Unmarshal(b, createLecturerDto); err != nil {
				return nil, err
			}

			// labId := sql.NullInt32 {}

			editor, err := queries.GetEditorByEmail(p.Context, editorEmail)

			if err != nil {
				fmt.Printf("failed on retrieving editor, %+v\n", err)
				return nil, gqlerrors.FormatError(errors.New("server side error"))
			}

			id, err := queries.CreateLecturer(p.Context, sqlcdb.CreateLecturerParams{
				Name:        createLecturerDto.Name,
				Email:       createLecturerDto.Email,
				Description: null.StringFromPtr(createLecturerDto.Description),
				Labid:       sql.NullInt64{Int64: int64(createLecturerDto.Labid), Valid: true},
				Gender:      null.StringFromPtr(createLecturerDto.Gender),
				Editorid:    sql.NullInt64{Int64: int64(editor.ID), Valid: true},
			})

			if err != nil {
				var mysqlErr *mysql.MySQLError
				if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
					return nil, errors.New("data already exists")
				}
				fmt.Printf("failed oh failed\n")
				return nil, err
			}
			data, _ := queries.GetLecturersByPk(p.Context, id)

			return data, nil
		},
	}
}

var UpdateLecturerInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "UpdateLecturersInput",
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
		"gender": &graphql.InputObjectFieldConfig{
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

type UpdateLecturerDto struct {
	Name        *string
	Email       *string
	Description *string
	Gender      *string
	Labid       *int
}

func existInt64(ptr *int) sql.NullInt64 {
	if ptr != nil {
		return sql.NullInt64{Int64: int64(*ptr), Valid: true}
	}
	return sql.NullInt64{}
}

func (factory *GqlFactory) UpdateLecturer(returnType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type: returnType,
		Args: graphql.FieldConfigArgument{
			"input": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(UpdateLecturerInput),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			inputMap, _ := p.Args["input"].(map[string]interface{})

			b, err := json.Marshal(inputMap)

			if err != nil {
				return nil, err
			}

			updateLecturerInput := &UpdateLecturerDto{}

			if err := json.Unmarshal(b, updateLecturerInput); err != nil {
				return nil, err
			}

			fmt.Printf("updateLecturerInput: %+v\n", updateLecturerInput)

			targetId, ok := inputMap["id"].(int)

			if !ok {
				return nil, gqlerrors.FormatError(errors.New("server side error"))
			}

			queries := sqlcdb.New(factory.DB)
			err = queries.UpdateLecturer(p.Context, sqlcdb.UpdateLecturerParams{
				ID:          int64(targetId),
				Name:        null.StringFromPtr(updateLecturerInput.Name),
				Email:       null.StringFromPtr(updateLecturerInput.Email),
				Description: null.StringFromPtr(updateLecturerInput.Description),
				Gender:      null.StringFromPtr(updateLecturerInput.Gender),
				Labid:       existInt64(updateLecturerInput.Labid),
			})

			if err != nil {
				return nil, err
			}

			var lecturer, _ = queries.GetLecturersByPk(p.Context, int64(targetId))

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
			fmt.Printf("----DeleteLecturersByPk was called-----\n")
			lecturerId, ok := p.Args["id"].(int)

			if !ok {
				return nil, gqlerrors.FormatError(errors.New("server side error"))
			}

			queries := sqlcdb.New(factory.DB)

			lecturer, err := queries.GetLecturersByPk(p.Context, int64(lecturerId))

			if err != nil {
				return nil, err
			}

			if err := queries.DeleteLecturersByPk(p.Context, int64(lecturerId)); err != nil {
				return nil, err
			}

			return lecturer, nil
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
				return nil, errors.New("invalid ids")
			}

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
