package db

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ondbyte/ogo/petstore_example/db/pets"
)

type Db struct {
	SqlDb *sql.DB
	Pets  *pets.Queries
}

func (q *Db) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return q.SqlDb.BeginTx(ctx, nil)
}

var connectionString = "myuser:yadu@tcp(localhost:3306)/mydatabase"

func InitDb(ctx context.Context, dbSchema string) (*Db, error) {
	var d, err = sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	_, err = d.ExecContext(ctx, dbSchema)
	if err != nil {
		return nil, err
	}
	return &Db{
		SqlDb: d,
		Pets:  pets.New(d),
	}, nil
}
