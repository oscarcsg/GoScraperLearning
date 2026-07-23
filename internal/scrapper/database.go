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
				price  BIGINT NOT NULL DEFAULT 0);`

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
	fmt.Fprint(&sb, "INSERT INTO books (title, rating, price) VALUES ")

	var args []any
	for i, book := range books.Books {
		sb.WriteString("(?, ?, ?)")

		// If it is not the last one, write a coma
		if i < len(books.Books)-1 {
			sb.WriteString(", ")
		} else {
			sb.WriteString(";")
		}

		args = append(
			args,
			strings.TrimSpace(book.Title),
			book.Rating,
			book.Price,
		)
	}
	query := sb.String()

	//fmt.Println(query)

	_, err := db.Exec(query, args...)
	if err != nil {
		logging.Error(
			"Books insertion went wrong.",
			logging.ErrorType(err),
		)
		return books
	} else {
		logging.Info(
			"Book page inserted in the database successfully.",
			logging.Uint16Type("web_page", books.WebPage),
			logging.IntType("books_count", len(books.Books)),
		)
	}
	
	return nil
}
