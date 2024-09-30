package postgres

import (
	"context"
	"log"
	"os"

	"github.com/np-d/boilerplate-go-fiber-sqlc/app/database/postgres/sqlc"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	Ctx        *context.Context
	Connection *pgx.Conn
	Queries    *sqlc.Queries
}

func Connect() *Database {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	return &Database{
		Ctx:        &ctx,
		Connection: conn,
		Queries:    sqlc.New(conn),
	}
}

func (db *Database) Close() {
	err := db.Connection.Close(*db.Ctx)
	if err != nil {
		log.Fatalf("unable to close database connection: %v\n", err)
	}
}
