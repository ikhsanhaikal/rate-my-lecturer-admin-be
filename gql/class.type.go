package gql

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
)

func (factory *GqlFactory) CourseType(subjectType *graphql.Object, lecturerType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "course",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.Int),
			},
			"lecturer": &graphql.Field{
				Type: graphql.NewNonNull(lecturerType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Class)
					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}
					queries := sqlcdb.New(factory.DB)

					if !source.Lecturerid.Valid {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					lecturer, err := queries.GetLecturersByPk(p.Context, source.Lecturerid.Int64)
					if err != nil {
						return nil, err
					}
					return lecturer, nil
				},
			},
			"subject": &graphql.Field{
				Type: graphql.NewNonNull(subjectType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Class)
					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}
					queries := sqlcdb.New(factory.DB)
					if !source.Subjectid.Valid {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					subject, err := queries.GetSubjectsByPk(p.Context, source.Subjectid.Int64)
					if err != nil {
						return nil, err
					}
					return subject, nil
				},
			},
			"semester": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Class)
					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					return source.Semester.String, nil
				},
			},
			"year": &graphql.Field{
				Type: graphql.DateTime,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Class)

					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					if source.Year.Valid {
						return source.Year.Time, nil
					}

					return nil, nil
				},
			},
			"createdAt": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	})
}
