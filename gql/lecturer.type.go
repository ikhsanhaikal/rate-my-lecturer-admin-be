package gql

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
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

					id := 0

					switch v := p.Source.(type) {
					case sqlcdb.Lecturer:
						fmt.Printf("lecturer type %v\n", v)
						source := p.Source.(sqlcdb.Lecturer)
						id = int(source.Labid)
					case sqlcdb.ListMembersRow:
						fmt.Printf("members type")
						source := p.Source.(sqlcdb.ListMembersRow)
						id = int(source.Labid)
					}

					queries := sqlcdb.New(builder.DB)

					lab, err := queries.GetLab(p.Context, int32(id))

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
