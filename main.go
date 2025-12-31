package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"tutorial.sqlc.dev/app/tutorial"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("set DATABASE_URL environment variable, e.g. postgres://user:pass@localhost:5432/dbname?sslmode=disable")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("opening db: %v", err)
	}
	defer db.Close()

	// Verify connection
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	q := tutorial.New(db)

	// Create a new author
	_, err = q.CreateAuthor(ctx, tutorial.CreateAuthorParams{
		Name: "Ada Lovelace",
		Bio:  sql.NullString{String: "First computer programmer", Valid: true},
	})
	if err != nil {
		log.Fatalf("CreateAuthor: %v", err)
	}

	// List authors
	authors, err := q.ListAuthors(ctx)
	if err != nil {
		log.Fatalf("ListAuthors: %v", err)
	}

	fmt.Println("Authors:")
	for _, a := range authors {
		bio := ""
		if a.Bio.Valid {
			bio = a.Bio.String
		}
		fmt.Printf("- %d: %s (%s)\n", a.ID, a.Name, bio)
	}
}
