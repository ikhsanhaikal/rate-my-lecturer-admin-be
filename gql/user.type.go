package gql

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
)

func (factory *GqlFactory) UserType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.User)

					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}
					return source.Username, nil
				},
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
}
