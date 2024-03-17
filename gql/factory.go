package gql

import "database/sql"

type GqlFactory struct {
	DB *sql.DB
}
