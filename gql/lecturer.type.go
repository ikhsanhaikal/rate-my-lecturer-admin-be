package gql

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
)

func (factory *GqlFactory) LecturerType(labType *graphql.Object) *graphql.Object {
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
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Lecturer)
					if !source.Description.Valid || !ok {
						return "", nil
					}
					return source.Description.String, nil
				},
			},
			"lab": &graphql.Field{
				Type: labType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {

					id := 0

					switch v := p.Source.(type) {
					case sqlcdb.Lecturer:
						fmt.Printf("lecturer type %+v\n", v)
						source := p.Source.(sqlcdb.Lecturer)
						id = int(source.Labid)
					case sqlcdb.ListMembersRow:
						fmt.Printf("members type %+v\n", v)
						source := p.Source.(sqlcdb.ListMembersRow)
						id = int(source.Labid)
					default:
						fmt.Printf("ohhh u screwed\n")
						fmt.Printf("type: %T\n", v)
					}

					queries := sqlcdb.New(factory.DB)

					lab, err := queries.GetLabsByPk(p.Context, int32(id))
					fmt.Printf("labid: %d\n", id)
					if err != nil {
						fmt.Printf("getLabsByPk err: %+v\n", err)
						return nil, err
					}

					fmt.Printf("lab: %+v\n", lab)

					return lab, nil
				},
			},
		},
	})
}
