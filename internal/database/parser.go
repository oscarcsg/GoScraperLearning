package database

import (
	"go-scraper-learning/internal/util"
)

func parseDriver(base string) (string) {
	// TODO: if I ensure the config importer to trim and lower correctly, this line must be gone
	base = util.TrimToLowerString(base)

	switch base {
	case "postgresql", "postgres":
		return "postgres"
	case "mysql", "mariadb":
		return "mysql"
	case "sqlite", "sqlite3":
		return "sqlite"

	default:
		return "unsupported"
	}
}