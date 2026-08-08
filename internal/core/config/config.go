package config

import "os"

type Config struct {
	DatabaseDSN string
}

func Load() Config {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://bop:bop@localhost:5432/bop?sslmode=disable"
	}

	return Config{DatabaseDSN: dsn}
}
