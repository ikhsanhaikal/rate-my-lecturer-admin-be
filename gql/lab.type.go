package gql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
)

func (factory *GqlFactory) LabType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "lab",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"code": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
}

func (r *Resolver) GetLecturersByLab(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)

	if !ok {
		return nil, errors.New("invalid pk")
	}
	queries := sqlcdb.New(r.DB)
	members, err := queries.ListMembers(p.Context, int32(id))

	if err != nil {
		return nil, err
	}

	fmt.Printf("lecturers_by_lab: %+v\n", members)

	return members, nil
}
func (r *Resolver) ListLabs(p graphql.ResolveParams) (interface{}, error) {
	queries := sqlcdb.New(r.DB)

	labs, err := queries.ListLabs(p.Context)

	fmt.Printf("labs: %+v\n", labs)

	if err != nil {
		return nil, err
	}

	return labs, nil
}
func (r *Resolver) GetLabById(p graphql.ResolveParams) (interface{}, error) {
	queries := sqlcdb.New(r.DB)

	id, ok := p.Args["id"].(int)

	if !ok {
		return nil, errors.New("server error")
	}

	lab, err := queries.GetLabsByPk(p.Context, int32(id))

	if err != nil {
		return nil, err
	}

	return lab, nil
}

func (factory *GqlFactory) CreateLab(returnType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        returnType,
		Description: "create a lab",
		Args: graphql.FieldConfigArgument{
			"input": &graphql.ArgumentConfig{
				Type: graphql.NewInputObject(graphql.InputObjectConfig{
					Name: "CreateLabInput",
					Fields: graphql.InputObjectConfigFieldMap{
						"name": &graphql.InputObjectFieldConfig{
							Type: graphql.String,
						},
						"code": &graphql.InputObjectFieldConfig{
							Type: graphql.String,
						},
						"description": &graphql.InputObjectFieldConfig{
							Type: graphql.String,
						},
					},
				}),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var queries = sqlcdb.New(factory.DB)

			input, ok := p.Args["input"].(map[string]interface{})

			if !ok {
				fmt.Printf("err can't cast args to map[string]interface{}\n")
				return nil, nil
			}

			fmt.Printf("input: %+v\n", input)

			result, err := queries.CreateLab(p.Context, sqlcdb.CreateLabParams{
				Name:        input["name"].(string),
				Code:        input["code"].(string),
				Description: sql.NullString{},
			})

			if err != nil {
				var mysqlErr *mysql.MySQLError
				if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
					return nil, errors.New("data already exists")
				}
				fmt.Printf("failed oh failed\n")
				return nil, err
			}

			id, _ := result.LastInsertId()

			data, _ := queries.GetLabsByPk(p.Context, int32(id))

			return data, nil
		},
	}
}

func (factory *GqlFactory) DeleteLabsByPk(returnType *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type:        returnType,
		Description: "delete lab",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var queries = sqlcdb.New(factory.DB)

			id, ok := p.Args["id"].(int)

			if !ok {
				return nil, nil
			}

			result, err := queries.DeleteLab(p.Context, int32(id))

			if err != nil {
				return nil, err
			}

			deletedId, _ := result.RowsAffected()

			return deletedId, nil
		},
	}
}
