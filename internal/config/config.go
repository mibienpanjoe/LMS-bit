package config

import "os"

type Config struct {
	AppName  string
	LogLevel string
}

func Load() Config {
	return Config{
		AppName:  getEnv("LMS_APP_NAME", "LMS-bit"),
		LogLevel: getEnv("LMS_LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	return v
}
