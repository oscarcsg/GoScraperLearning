module go-scraper-learning

go 1.26.1

// Extraer variables de entorno
require github.com/joho/godotenv v1.5.1

// Goquery (web scrapping)
require (
	github.com/PuerkitoBio/goquery v1.12.0 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	golang.org/x/net v0.52.0 // indirect
)

// Databases drivers
require (
	// PostgreSQL
	github.com/jackc/pgx v3.6.2+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/text v0.35.0 // indirect

	// SQLite
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/sys v0.44.0 // indirect
	modernc.org/libc v1.73.4 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/sqlite v1.53.0 // indirect

	// MySQL y MariaDB
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/go-sql-driver/mysql v1.10.0 // indirect
)

require (
	
)
