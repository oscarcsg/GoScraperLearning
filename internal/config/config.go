package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Database
func loadDatabaseConfig() (DatabaseConfig) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("IMPORTANTE: archivo '.env' no encontrado. Leyendo las variables del sistema.")
	}
	
	engine := strings.ToLower(os.Getenv("DB_ENGINE"))
	if engine == "" {
		fmt.Println("FATAL: el motor de la base de datos DEBE estar definido.")
		os.Exit(1)
	}

	switch engine {
	case "sqlite":
		fmt.Println("Se usará una base de datos SQLite.")
	case "postgresql":
		fmt.Println("Se usará una base de datos PostgreSQL.")
	case "mysql":
		fmt.Println("Se usará una base de datos SQLite.")
	case "sqlserver":
		fmt.Println("Se usará una base de datos SQLite.")
	case "mariadb":
		fmt.Println("Se usará una base de datos SQLite.")
	
	default:
		fmt.Println("Opción no válida.")
	}

	return DatabaseConfig {
		DBEngine: engine,
		DBSQLite: "",
		DBUser: "",
		DBUserPassword: "",
		DBIpUrl: "",
		DBPort: "",
		DBName: "",
	}
}



// ------ GLOBAL PUBLIC METHOD ------ //
func LoadConfig() (DatabaseConfig) {
	return loadDatabaseConfig()
}