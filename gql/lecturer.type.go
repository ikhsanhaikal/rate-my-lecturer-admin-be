package gql

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
	db "github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/mysql"
)

func (builder TypeBuilder) LecturerType(labType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "lecturer",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"lab": &graphql.Field{
				Type: labType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(db.Lecturer)

					if !ok {
						return nil, errors.New("failed that pretty much u need to know :).")
					}

					queries := db.New(builder.DB)

					// lab, err := source.getLab()

					lab, err := queries.GetLab(p.Context, int32(source.Labid))

					if err != nil {
						return nil, err
					}

					fmt.Printf("lab: %+v\n", lab)

					return lab, nil
				},
			},
		},
	})
}
