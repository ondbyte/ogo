package petstoreexample

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ondbyte/ogo/petstore_example/db/pets"
)

func initRepos(ctx context.Context, db *sql.DB) (*pets.Queries, error) {
	return pets.New(db), nil
}
