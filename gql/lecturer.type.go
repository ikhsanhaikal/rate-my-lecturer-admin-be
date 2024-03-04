package gql

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
	db "github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/mysql"
)

func (r Repository) LecturerType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Lecturer",
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
					fmt.Printf("Resolve description: i was called\n")
					source, ok := p.Source.(db.Lecturer)

					if !ok {
						return nil, errors.New("mehh")
					}

					fmt.Printf("source: %+v\n", source)

					return source.Description, nil
				},
			},
			"lab": &graphql.Field{
				Type: LabType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(db.Lecturer)

					if !ok {
						return nil, errors.New("mehh")
					}

					queries := db.New(r.Db)

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

// var LecturerType = graphql.NewObject(graphql.ObjectConfig{
// 	Name: "Lecturer",
// 	Fields: graphql.Fields{
// 		"id": &graphql.Field{
// 			Type: graphql.String,
// 		},
// 		"name": &graphql.Field{
// 			Type: graphql.String,
// 		},
// 		"email": &graphql.Field{
// 			Type: graphql.String,
// 		},
// 		"description": &graphql.Field{
// 			Type: graphql.String,
// 			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 				fmt.Printf("Resolve description: i was called\n")
// 				return "laskdmakl", nil
// 			},
// 		},
// 		"field1": &graphql.Field{
// 			Type: graphql.Boolean,
// 			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 				fmt.Printf("\nfield1: alksmaks\n")
// 				return false, nil
// 			},
// 		},
// 	},
// })
