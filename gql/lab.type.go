package gql

import (
	"github.com/graphql-go/graphql"
)

var LabType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Lab",
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
