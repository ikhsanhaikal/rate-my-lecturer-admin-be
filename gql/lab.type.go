package gql

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
)

func (builder TypeBuilder) LabType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "lab",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"code": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
}

func (r *Resolver) GetLecturersByLab(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)

	if !ok {
		return nil, errors.New("invalid pk")
	}
	queries := sqlcdb.New(r.DB)
	members, err := queries.ListMembers(p.Context, int32(id))

	if err != nil {
		return nil, err
	}

	fmt.Printf("lecturers_by_lab: %+v\n", members)

	return members, nil
}

func (r *Resolver) ListLabs(p graphql.ResolveParams) (interface{}, error) {
	queries := sqlcdb.New(r.DB)

	labs, err := queries.ListLabs(p.Context)

	fmt.Printf("labs: %+v\n", labs)

	if err != nil {
		return nil, err
	}

	return labs, nil
}
func (r *Resolver) GetLabById(p graphql.ResolveParams) (interface{}, error) {
	queries := sqlcdb.New(r.DB)

	id, ok := p.Args["id"].(int)

	if !ok {
		return nil, errors.New("server error")
	}

	lab, err := queries.GetLab(p.Context, int32(id))

	if err != nil {
		return nil, err
	}

	return lab, nil
}
