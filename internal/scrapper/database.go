package scrapper

import (
	"database/sql"
	"go-scraper-learning/internal/logging"
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

func InsertBooks(books []BookRegisterDTO) {
	query := ``
}