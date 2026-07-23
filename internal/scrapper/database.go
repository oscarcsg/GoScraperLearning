package scrapper

import (
	"database/sql"
	"fmt"
	"go-scraper-learning/internal/logging"
	"strings"
)

func CreateDatabaseSQLite(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS books(
			  	id     INTEGER PRIMARY KEY AUTOINCREMENT,
			  	title  TEXT NOT NULL,
				rating TINYINT NOT NULL DEFAULT 0,
				price  DOUBLE NOT NULL DEFAULT 0.0);`

	rows, err := db.Query(query)
	if err != nil {
		logging.Fatal(
			"SQLite database table creation failed.",
			logging.ErrorType(err),
		)
		return
	}
	defer rows.Close()

	logging.Info(
		"Table created successfully in the SQLite database.",
	)
}

func InsertBooks(db *sql.DB, books *BooksPage) (*BooksPage) {
	var sb strings.Builder
	fmt.Fprint(&sb, "INSERT INTO books (title, rating, price) ")
	//query := "INSERT INTO books (title, rating, price) "
	for _, book := range books.Books {
		fmt.Fprintf(
			&sb,
			"VALUES(%s, %d, %d),",
			strings.TrimSpace(book.Title),
			book.Rating,
			book.Price,
		)
	}
	query := sb.String()
	query = strings.TrimSuffix(query, ",")
	query = query + ";"

	_, err := db.Exec(query)
	if err != nil {
		logging.Error(
			"Books insertion went wrong.",
			logging.ErrorType(err),
		)
		return books
	}
	
	return nil
}