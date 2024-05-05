package gql

import "github.com/graphql-go/graphql"

func (factory *GqlFactory) TraitType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "trait",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
}
