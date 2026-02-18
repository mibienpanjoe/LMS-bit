package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppName         string
	LogLevel        string
	StoragePath     string
	LoanDays        int
	MaxLoansPerUser int
	MaxLoanRenewals int
}

func Load() Config {
	return Config{
		AppName:         getEnv("LMS_APP_NAME", "Library Management System"),
		LogLevel:        getEnv("LMS_LOG_LEVEL", "info"),
		StoragePath:     getEnv("LMS_STORAGE_PATH", "data/storage.json"),
		LoanDays:        getEnvInt("LMS_LOAN_DAYS", 14),
		MaxLoansPerUser: getEnvInt("LMS_MAX_LOANS_PER_MEMBER", 3),
		MaxLoanRenewals: getEnvInt("LMS_MAX_LOAN_RENEWALS", 1),
	}
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	return v
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}

	return n
}
