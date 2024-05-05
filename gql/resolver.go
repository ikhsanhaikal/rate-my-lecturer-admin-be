package gql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/graphql-go/graphql"
	"github.com/ikhsanhaikal/rate-my-lecturer-graphql-admin/be-app/sqlcdb"
)

type Resolver struct {
	DB *sql.DB
}
type ListLecturersResult struct {
	Data  []sqlcdb.Lecturer
	Total int64
}

func (r *Resolver) ListLecturers(p graphql.ResolveParams) (interface{}, error) {
	queries := sqlcdb.New(r.DB)

	limit, _ := p.Args["limit"].(int)
	page, _ := p.Args["page"].(int)
	sortMap, _ := p.Args["sort"].(map[string]interface{})

	// fmt.Printf("limit: %d, page: %d\n", limit, page)

	if sortMap != nil {
		fmt.Printf("\nnot nil sort: %+v\n", sortMap)
		bites, _ := json.Marshal(sortMap)
		sort := &struct {
			Field string
			Order string
		}{}
		json.Unmarshal(bites, sort)
		if sort.Order == "DESC" {
			lecturers, err := queries.ListLecturersDesc(p.Context, sqlcdb.ListLecturersDescParams{
				Limit:  int32(limit),
				Offset: int32(limit) * int32(page-1),
			})
			if err != nil {
				return nil, err
			}
			total, err := queries.CountLecturers(p.Context)

			if err != nil {
				return nil, err
			}

			return ListLecturersResult{
				Data:  lecturers,
				Total: total,
			}, nil

		}
	} else {
		fmt.Printf("\nnil sort: %+v\n", sortMap)
	}

	lecturers, err := queries.ListLecturers(p.Context, sqlcdb.ListLecturersParams{
		Limit:  int32(limit),
		Offset: int32(limit) * int32(page-1),
	})
	if err != nil {
		return nil, err
	}
	total, err := queries.CountLecturers(p.Context)

	if err != nil {
		return nil, err
	}

	return ListLecturersResult{
		Data:  lecturers,
		Total: total,
	}, nil
}

func (r *Resolver) GetLecturerById(p graphql.ResolveParams) (interface{}, error) {
	queries := sqlcdb.New(r.DB)

	id := p.Args["id"].(int)

	lecturer, err := queries.GetLecturersByPk(p.Context, int64(id))

	if err != nil {
		return nil, err
	}

	return lecturer, nil
}

func (r *Resolver) CreateUser(p graphql.ResolveParams) (interface{}, error) {
	var queries = sqlcdb.New(r.DB)

	input, ok := p.Args["input"].(map[string]interface{})

	if !ok {
		fmt.Printf("err can't cast args to struct\n")
		return nil, nil
	}

	result, err := queries.CreateUser(p.Context, sqlcdb.CreateUserParams{
		Username: input["name"].(string),
		Email:    input["email"].(string),
	})

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, errors.New("data already exists")
		}
		fmt.Printf("failed oh failed\n")
		return nil, err
	}

	id, _ := result.LastInsertId()

	data, _ := queries.GetUser(p.Context, int64(id))

	return data, nil
}
