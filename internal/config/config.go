package config

import (
	"fmt"
	"go-scraper-learning/internal/util"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// ---------------------------------- //
// -------- DATABASE METHODS -------- //
// ---------------------------------- //
// (return: engine, conn string and error)
func loadDatabaseConfig() (string, string, error) {
	engine := util.TrimToLowerString(os.Getenv("DB_ENGINE"))
	//fmt.Println(engine)
	if engine == "" {
		fmt.Println("FATAL: database engine MUST be defined.") // Log
		os.Exit(1)
	}

	switch engine {
	case "sqlite":
		dbPath := strings.TrimSpace(os.Getenv("DB_SQLITE"))
		if dbPath == "" {
			return "", "", fmt.Errorf("FATAL: SQLite database path MUST be defined.")
		}
		return engine,
			   fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)", dbPath),
			   nil

	case "postgresql", "mysql", "mariadb":
		// Get the values from the env (and possible errors)
		values, errors := getAndCheckServerDatabaseParameters()

		if len(errors) > 0 {
			for _, err := range errors {
				fmt.Println(err) // Log
			}
			os.Exit(1)
		}

		// Values should be in this order: user, userPsw, ipUrl, port, dbName
		if len(values) != 5 {
			return "", "", fmt.Errorf("FATAL: missing database server connection values")
		}

		if engine == "postgresql" {
			return engine,
				fmt.Sprintf(
					"postgres://%s:%s@%s:%s/%s?sslmode=disable",
					values[0],
					values[1],
					values[2],
					values[3],
					values[4],
				),
				nil
		}

		return engine,
			fmt.Sprintf(
				"%s:%s@tcp(%s:%s)/%s?parseTime=true",
				values[0],
				values[1],
				values[2],
				values[3],
				values[4],
			),
			nil

	default:
		return "", "", fmt.Errorf("Invalid option: %s.", engine)
	}
}

func getAndCheckServerDatabaseParameters() ([]string, []error) {
	// Lenght 0 (how many objects already have), Capacity 5 (the max num of objects it can hold)
	values := make([]string, 0, 5)
	errors := make([]error, 0, 5)

	errorMsgs := map[string]string{
		"DB_USER":          "Database server user not found in the environment.",
		"DB_USER_PASSWORD": "Database server user password not found in the environment.",
		"DB_IP_URL":        "Database server IP or URL not found in the environment.",
		"DB_PORT":          "Database server port not found in the environment.",
		"DB_NAME":          "Database name not found in the environment.",
	}

	// Keys in the .env
	keys := []string{"DB_USER", "DB_USER_PASSWORD", "DB_IP_URL", "DB_PORT", "DB_NAME"}

	// Iterate over every key
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			errors = append(errors, fmt.Errorf("FATAL: %s", errorMsgs[key]))
		} else {
			values = append(values, value)
		}
	}

	return values, errors
}


// ---------------------------------- //
// -------- LOGGING METHODS --------- //
// ---------------------------------- //
// (return: struct LogConfig)
func loadLogConfig() (LogConfig) {
	errors := make([]error, 0, 6)

	config := LogConfig{
		FilePath: getStringEnv("LOG_FILE_NAME", &errors),
		FileMaxSize: getUint16Env("LOG_FILE_MAX_SIZE", &errors),
		FileMaxAge: getUint16Env("LOG_FILE_MAX_AGE", &errors),
		FileMaxBackups: getUint8Env("LOG_FILE_MAX_BACKUPS", &errors),
		FileCompress: getBoolEnv("LOG_FILE_COMPRESS", &errors),
		Terminal: getBoolEnv("LOG_TERMINAL", &errors),
	}

	// Print errors and exit the program
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err) // Log
		}
		os.Exit(1)
	}

	// Create the logconfig and return
	return config
}



// ---------------------------------- //
// -------- TEL LOG METHODS --------- //
// ---------------------------------- //
// (return: struct TelegramLogConfig)




// ---------------------------------- //
// ------ GLOBAL PUBLIC METHOD ------ //
// ---------------------------------- //
func LoadConfig()  {
	if err := godotenv.Load(); err != nil {
		fmt.Println("IMPORTANTE: '.env' file not found. Reading system environment variables...")
	}

	dbEngine, dbString, err := loadDatabaseConfig()

	if err != nil {
		fmt.Println(err) // Log
		os.Exit(1)
	}

	fmt.Printf("\nEngine: %s; Connection String: %s", dbEngine, dbString)

	logConfig := loadLogConfig()

	loadTelegramLogConfig()
}



// ---------------------------------- //
// -------- UTILITY METHODS --------- //
// ---------------------------------- //
func getStringEnv(envKey string, errors *[]error) (string) {
	txt := strings.TrimSpace(os.Getenv(envKey))
	if txt == "" {
		*errors = append(*errors, fmt.Errorf("FATAL: %s variable is empty or does not exist.", envKey))
	}
	return txt
}

func getStringEnvTrimLower(envKey string, errors *[]error) (string) {
	txt := util.TrimToLowerString(os.Getenv(envKey))
	if txt == "" {
		*errors = append(*errors, fmt.Errorf("FATAL: %s variable is empty or does not exist.", envKey))
	}
	return txt
}

func getUint8Env(envKey string, errors *[]error) (uint8) {
	txt := getStringEnv(envKey, errors)
	if txt == "" {
		// Error msg has been already written in the previous method
		return 0
	}

	value, err := util.ParseUint8(txt)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("FATAL: %s variable value is not valid.", envKey))
		return 0
	}
	return value
}

func getUint16Env(envKey string, errors *[]error) (uint16) {
	txt := getStringEnv(envKey, errors)
	if txt == "" {
		// Error msg has been already written in the previous method
		return 0
	}

	value, err := util.ParseUint16(txt)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("FATAL: %s variable value is not valid.", envKey))
		return 0
	}
	return value
}

func getBoolEnv(envKey string, errors *[]error) (bool) {
	txt := getStringEnv(envKey, errors)
	if txt == "" {
		// Error msg has been already written in the previous method
		return false
	}

	value, err := strconv.ParseBool(txt)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("FATAL: %s variable value is not valid.", envKey))
		return false
	}
	return value
}
