package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gen1usBruh/url_shortener/internal/config"
	"github.com/Gen1usBruh/url_shortener/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

const (
	unique_violation string = "23505"
)

/*
Notes: Database, schema, tables
were already created.

CREATE TABLE urlschema.urls (

	id SERIAL PRIMARY KEY,
	alias TEXT NOT NULL UNIQUE,
	url TEXT NOT NULL

);

CREATE INDEX idx_alias ON urlschema.urls(alias);
*/
func ConnectDB(dbConfig *config.Database) (*Storage, error) {

	const op = "storage.postgres.ConnectDB"

	// urlExample := "postgres://username:password@localhost:port/database_name"
	dbUrl := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		dbConfig.User, dbConfig.Password, dbConfig.Host,
		dbConfig.Port, dbConfig.DBname,
	)

	dbPool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s | %w", op, err)
	}

	if err = dbPool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s | %w", op, err)
	}

	return &Storage{db: dbPool}, nil
}

/*
Note: we use QueryRow() instead of Exec() since
PostgreSQL does not return you the last inserted id.
This is because last inserted id is available only if
you create a new row in a table that uses a sequence.

If you actually insert a row in the table where
a sequence is assigned, you have to use RETURNING clause.
*/
func (s *Storage) SaveURL(urlToSave string, alias string) (int32, error) {
	const op = "storage.postgres.SaveURL"

	var id int32
	stmt := `INSERT INTO urlschema.urls(url, alias) VALUES ($1, $2) RETURNING id`
	err := s.db.QueryRow(context.Background(), stmt, urlToSave, alias).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == unique_violation {
			return 0, fmt.Errorf("%s, %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s | %w", op, err)
	}

	return id, nil
}
