package gql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
)

func (factory *GqlFactory) LecturerType(labType, traitType *graphql.Object) *graphql.Object {
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
			"gender": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Lecturer)

					if !ok {
						log.Printf("LecturerType gender field resolve error\n")
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					if !source.Gender.Valid {
						return "", nil
					}

					return source.Gender.String, nil
				},
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

					// id := 0

					source, _ := p.Source.(sqlcdb.Lecturer)
					// switch v := p.Source.(type) {
					// case sqlcdb.Lecturer:
					// source := p.Source.(sqlcdb.Lecturer)
					// id = int(source.Labid)
					// case sqlcdb.ListMembersRow:
					// 	fmt.Printf("members type %+v\n", v)
					// 	source := p.Source.(sqlcdb.ListMembersRow)
					// 	id = int(source.Labid)
					// default:
					// fmt.Printf("ohhh u screwed\n")
					// fmt.Printf("type: %T\n", v)
					// }

					queries := sqlcdb.New(factory.DB)

					if !source.Labid.Valid {
						return nil, nil
					}

					lab, err := queries.GetLabsByPk(p.Context, source.Labid.Int64)

					if err != nil {
						log.Printf("LecturerType lab field resolve error\n")
						return nil, err
					}

					return lab, nil
				},
			},
			"rating": &graphql.Field{
				Type: graphql.Float,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Lecturer)

					if !ok {
						return nil, graphql.NewLocatedError("failed on type assertion source ", nil)
					}

					queries := sqlcdb.New(factory.DB)
					result, err := queries.AverageRatingByLecturerId(context.Background(), sql.NullInt64{Int64: int64(source.ID), Valid: true})

					if err != nil {
						return nil, err
					}

					avg, ok := result.(float64)

					if !ok {
						fmt.Printf("\nNot Ok type casting\n")
						return 0, nil
					}

					return avg, nil
				},
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Lecturer)
					if !ok {
						log.Printf("LecturerType created_at field resolve error\n")
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}
					return source.Createdat, nil
				},
			},
			"tags": &graphql.Field{
				Type: graphql.NewList(traitType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Lecturer)

					if !ok {
						log.Printf("LecturerType tags field resolve error\n")
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					queries := sqlcdb.New(factory.DB)

					tags, err := queries.GetSummaryTagsByLecturer(p.Context, sql.NullInt64{Int64: int64(source.ID), Valid: true})

					if err != nil {
						return nil, err
					}

					return tags, nil
				},
			},
			"editor": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Lecturer)

					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					queries := sqlcdb.New(factory.DB)
					editor, err := queries.GetEditor(p.Context, int64(source.Editorid.Int64))

					if err != nil {
						return nil, err
					}

					return editor.Email, nil
				},
			},
		},
	})
}

func (factory *GqlFactory) ListLecturersResult(lecturerType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "ListLecturersResult",
		Fields: graphql.Fields{
			"data": &graphql.Field{
				Type: graphql.NewList(lecturerType),
			},
			"total": &graphql.Field{
				Description: "Total all lecturers in table",
				Type:        graphql.Int,
			},
		},
	})
}
