package gql

import "github.com/graphql-go/graphql"

var CreateLecturersInput = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "CreateLecturersInput",
	Fields: graphql.InputObjectConfigFieldMap{
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
		"gender": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})
