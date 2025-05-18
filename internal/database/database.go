package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/detectivekaktus/JGame/internal/config"
	"github.com/jackc/pgx/v5"
)

// Establishes the connection to the PostgreSQL database.
// You must first set up DATABASE_URL environment variable in
// the `.env` file inside the root directory of the project.
// Read README.md for more information.
//
// Don't forget to close the connection on the database once
// you're done using it.
func getConnection() *pgx.Conn {
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
func QueryRow(query string, args ...any) pgx.Row {
	conn := getConnection()
	defer conn.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	return conn.QueryRow(ctx, query, args...)
}

// Queries rows from the database as backend user.
// This function wraps the pgx.Conn.Query for
// simplicity and consistency.
// Don't forget to close the rows after you're done
// using them.
//
// TODO: Introduce parameter timeoutSec for handling the
// context timeout time for the query.
func QueryRows(query string, args ...any) pgx.Rows {
	conn := getConnection()
	defer conn.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not query %s: %v\n", query, err)
	}
	return rows
}

// Executes an SQL statement as backend user on the database.
// Return 0 on success, 1 on failure.
func Execute(stmnt string, args ...any) int {
	conn := getConnection()
	defer conn.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	_, err := conn.Exec(ctx, stmnt, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not exec %s: %v\n", stmnt, err)
		return 1
	}
	return 0
}
