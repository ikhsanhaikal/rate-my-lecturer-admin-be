package gql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
)

func (factory *GqlFactory) ReviewType(courseType *graphql.Object, traitType *graphql.Object, userType *graphql.Object) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "review",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"reviewer": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Review)

					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					queries := sqlcdb.New(factory.DB)
					user, err := queries.GetUser(p.Context, source.Reviewerid)

					if err != nil {
						return nil, err
					}

					return user, nil
				},
			},
			"comment": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Review)
					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}
					return source.Comment.String, nil
				},
			},
			"rating": &graphql.Field{
				Type: graphql.Float,
			},
			"course": &graphql.Field{
				Type: courseType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Review)

					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					queries := sqlcdb.New(factory.DB)
					course, err := queries.GetCourseByReview(p.Context, source.ID)

					if err != nil {
						return nil, err
					}

					return course, nil
				},
			},
			"tags": &graphql.Field{
				Type: graphql.NewList(traitType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source, ok := p.Source.(sqlcdb.Review)

					if !ok {
						return nil, gqlerrors.FormatError(errors.New("server side error"))
					}

					queries := sqlcdb.New(factory.DB)
					traits, err := queries.GetTagsByReview(p.Context, source.ID)

					if err != nil {
						return nil, err
					}

					fmt.Printf("tag resolver() traits: %+v\n", traits)
					return traits, nil
				},
			},
		},
	})
}

func (factory *GqlFactory) GetReviewsByLecturer(p graphql.ResolveParams) (interface{}, error) {
	lecturerId, ok := p.Args["id"].(int)
	page, _ := p.Args["page"].(int)
	limit, limitOk := p.Args["limit"].(int)

	fmt.Printf("GetReviewsByLecturer: (page: %+v), (limit: %+v)\n", page, limit)
	if !ok {
		log.Printf("Unable type assertion id: lecturerId\n")
		return nil, gqlerrors.FormatError(errors.New("server side error"))
	}

	if !limitOk {
		log.Printf("Unable type assertion limit\n")
		log.Printf("Unable type assertion limit: %+v\n", limit)
		return nil, gqlerrors.FormatError(errors.New("server side error"))
	}

	queries := sqlcdb.New(factory.DB)

	fmt.Printf("limit: %+v\n", limit)
	fmt.Printf("limit: %+v\n", int32(limit))
	reviews, err := queries.GetReviewsByLecturer(p.Context, sqlcdb.GetReviewsByLecturerParams{
		Lecturerid: sql.NullInt64{Int64: int64(lecturerId), Valid: true},
		Limit:      int32(limit),
		Offset:     int32(limit) * int32(page-1),
	})

	if err != nil {
		fmt.Printf("Err review type resolver 2: %+v\n", err)
		return nil, err
	}

	total, err := queries.CountReviewsByLecturer(p.Context, sql.NullInt64{Int64: int64(lecturerId), Valid: true})

	if err != nil {
		fmt.Printf("Err review type resolver 2: %+v\n", err)
		return nil, err
	}

	//debug only
	// if len(reviews) > 0 {
	// 	review := reviews[1]
	// 	bites, err := json.Marshal(review)
	// 	if err != nil {
	// 		fmt.Printf("\n\n\nfailed on marshal review %+v\n\n\n", err)
	// 	} else {
	// 		fmt.Printf("\n\n\nbites: %+v\n\n\n", string(bites))
	// 	}
	// }

	return struct {
		Data  []sqlcdb.Review
		Total int64
	}{
		Data:  reviews,
		Total: total,
	}, nil
}
