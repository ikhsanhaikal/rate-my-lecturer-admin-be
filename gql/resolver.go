package gql

import (
	"database/sql"
	"fmt"

	"github.com/graphql-go/graphql"
	db "github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/mysql"
)

type Resolver struct {
	DB *sql.DB
}

type TypeBuilder struct {
	DB *sql.DB
}

func (r *Resolver) ListLecturers(p graphql.ResolveParams) (interface{}, error) {
	queries := db.New(r.DB)

	lecturers, err := queries.ListLecturers(p.Context)

	if err != nil {
		return nil, err
	}

	return lecturers, nil
}

func (r *Resolver) GetLecturerById(p graphql.ResolveParams) (interface{}, error) {
	queries := db.New(r.DB)

	id := p.Args["id"].(int)

	lecturer, err := queries.GetLecturerById(p.Context, int32(id))

	fmt.Printf("lecturer: %+v\n", lecturer)

	if err != nil {
		return nil, err
	}

	return lecturer, nil
}
