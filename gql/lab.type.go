package gql

import (
	"github.com/graphql-go/graphql"
)

func (builder TypeBuilder) LabType() *graphql.Object {
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
		},
	})
}
