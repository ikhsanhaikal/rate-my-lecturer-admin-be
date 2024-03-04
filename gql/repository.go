package gql

import (
	"database/sql"

	"github.com/graphql-go/graphql"
	db "github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/mysql"
)

type Repository struct {
	Db *sql.DB
}

func (r Repository) ListLecturersResolver(p graphql.ResolveParams) (interface{}, error) {
	queries := db.New(r.Db)

	lecturers, err := queries.ListLecturers(p.Context)

	if err != nil {
		return nil, err
	}

	return lecturers, nil
}
