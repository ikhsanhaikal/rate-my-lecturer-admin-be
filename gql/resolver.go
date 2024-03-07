package gql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
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

func (r *Resolver) CreateLecturer(p graphql.ResolveParams) (interface{}, error) {
	var queries = db.New(r.DB)

	input, ok := p.Args["input"].(map[string]interface{})

	if !ok {
		fmt.Printf("err can't cast args to struct\n")
		return nil, nil
	}

	fmt.Printf("input: %+v\n", input)

	result, err := queries.CreateLecturer(p.Context, db.CreateLecturerParams{
		Name:        input["name"].(string),
		Email:       input["email"].(string),
		Description: input["description"].(string),
		Labid:       int32(input["labId"].(int)),
		Gender:      db.LecturersGenderMale,
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

	data, _ := queries.GetLecturerById(p.Context, int32(id))

	return data, nil
}

func (r *Resolver) CreateUser(p graphql.ResolveParams) (interface{}, error) {
	var queries = db.New(r.DB)

	input, ok := p.Args["input"].(map[string]interface{})

	if !ok {
		fmt.Printf("err can't cast args to struct\n")
		return nil, nil
	}

	fmt.Printf("input: %+v\n", input)

	result, err := queries.CreateUser(p.Context, db.CreateUserParams{
		Name:  input["name"].(string),
		Email: input["email"].(string),
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
