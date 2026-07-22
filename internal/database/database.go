package database

import (
	"database/sql"
	"fmt"
	"go-scraper-learning/internal/logging"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx"
	_ "modernc.org/sqlite"
)

func Init(config DatabaseConfig) (*sql.DB, error) {
	engineRaw, connString, errBool := getEngineAndConnString(config)
	// err is a bool
	if errBool { return nil, nil }

	engine := parseDriver(engineRaw)
	if engine == "unsupported" {
		logging.Error(
			"Database engine not supported.",
			logging.StringType("raw", engineRaw),
		)
	}

	isSQLite := false
	if engine == "sqlite" { isSQLite = true }

	if isSQLite {
		// Check or create the directory of the SQLite database file
		createSQLiteDBDirectory(connString)
	}

	db, err := sql.Open(engine, connString)
	if err != nil {
		logging.Error(
			"sql.Open() function error.",
			logging.ErrorType(err),
		)
		return nil, err
	}

	if isSQLite {
		// In sqlite only 1 connection is perfect, and no lifetime to avoid
		// closing and opening the file very X minutes
		db.SetMaxOpenConns(1)
		db.SetConnMaxLifetime(0)

		db.SetMaxIdleConns(1)
		db.SetConnMaxIdleTime(0)
	} else {
		db.SetMaxOpenConns(25)
		db.SetConnMaxLifetime(5 * time.Minute) // 5 minutes

		db.SetMaxIdleConns(25)
		db.SetConnMaxIdleTime(5 * time.Minute)
	}

	err = db.Ping()
	if err != nil {
		logging.Error(
			"Failed while ping. Database could not be initialized.",
			logging.ErrorType(err),
		)
		return nil, err
	}

	logging.Info(
		"Database connection established successfully.",
		logging.StringType("engine", engineRaw),
	)

	return db, nil
}

// Returns: engine, conn string
func getEngineAndConnString(config DatabaseConfig) (string, string, bool) {
	errs := make([]error, 0, 2)

	engine := strings.TrimSpace(config.DBEngine)
	if engine == "" {
		errs = append(errs, fmt.Errorf("Database engine not defined. Given engine: %s", engine))
	}

	connString := strings.TrimSpace(config.DBConnString)
	if connString == "" {
		logging.Error(
			"Database connection string not defined.",
			logging.StringType("connection_string", connString),
		)
		errs = append(errs, fmt.Errorf("Database connection string not defined. Given connection string: %s", connString))
	}

	if len(errs) > 0 {
		return engine, connString, true
	}
	return engine, connString, false
}

func createSQLiteDBDirectory(connString string) (bool) {
	rootDirectory := strings.HasPrefix(connString, "/")
	directories := strings.Split(strings.Split(connString, "?")[0], "/")

	var sb strings.Builder

	if rootDirectory {
		fmt.Fprint(&sb, "/")
	}

	for _, directory := range directories {
		if directory != "" && !strings.Contains(directory, ".db") {
			if strings.Contains(directory, "file:") {
				directory = strings.ReplaceAll(directory, "file:", "", )
			}
			// Valid directory
			fmt.Fprintf(&sb, "%s/", directory)
		}
	}

	err := os.MkdirAll(sb.String(), 0755) // drwxr-xr-x
	if err != nil {
		logging.Error(
			"Could not create the SQLite database directory. It will probably fail to connect.",
			logging.StringType("path", sb.String()),
		)
		return false
	}
	return true
}
