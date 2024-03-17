package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
)

func (factory *GqlFactory) SubjectType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "subject",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Subject)
					if !source.Description.Valid || !ok {
						return "", nil
					}
					return source.Description.String, nil
				},
			},
		},
	})
}
