package database

type FullDatabaseConfig struct {
	DBEngine       string
	DBSQLite       string
	DBUser         string
	DBUserPassword string
	DBIpUrl        string
	DBPort         string
	DBName         string
}

type DatabaseConfig struct {
	DBEngine     string
	DBConnString string
}