package config

import "os"

type Config struct {
	AppPort string
	DBDsn   string
}

func Load() Config {
	return Config{
		AppPort: os.Getenv("APP_PORT"),
		DBDsn:   os.Getenv("DATABASE_URL"),
	}
}
