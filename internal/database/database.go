package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/detectivekaktus/JGame/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	UniqueViolation = "23505"
	ForeignKeyViolation = "23503"
)

// Establishes the connection to the PostgreSQL database.
// You must first set up DATABASE_URL environment variable in
// the `.env` file inside the root directory of the project.
// Read README.md for more information.
//
// Don't forget to close the connection on the database once
// you're done using it.
func GetConnection() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), config.AppConfig.DbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to the database: %v\n", err)
		os.Exit(1)
	}
	return conn;
}

// Queries one row from the database as backend user.
// This function wraps the pgx.Conn.QueryRow for
// simplicity and consistency.
//
// TODO: Introduce parameter timeoutSec for handling the
// context timeout time for the query.
func QueryRow(conn *pgx.Conn, query string, args ...any) pgx.Row {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	return conn.QueryRow(ctx, query, args...)
}

// Queries rows from the database as backend user.
// This function wraps the pgx.Conn.Query for
// simplicity and consistency.
//
// Don't forget to close the rows after you're done
// using them.
//
// TODO: Introduce parameter timeoutSec for handling the
// context timeout time for the query.
func QueryRows(conn *pgx.Conn, query string, args ...any) pgx.Rows {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not query %s: %v\n", query, err)
	}
	return rows
}

// Executes an SQL statement as backend user on the database.
// The function wraps around pgx.Conn.Exec, see it for the return
// values.
func Execute(conn *pgx.Conn, stmnt string, args ...any) (pgconn.CommandTag, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	retVal, err := conn.Exec(ctx, stmnt, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not exec %s: %v\n", stmnt, err)
	}
	return retVal, err
}
