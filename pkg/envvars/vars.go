package env_vars

import (
	"os"
	"strconv"
)

const (
	IS_DEV_KEY                         = "IS_DEV"
	DEBUG_MODE_KEY                     = "DEBUG_MODE"
	GOOGLE_APPLICATION_CREDENTIALS_KEY = "GOOGLE_APPLICATION_CREDENTIALS"
	GOOGLE_PROJECT_ID_KEY              = "GOOGLE_PROJECT_ID_"

	IS_DEV_VALUE   = "1"
	IS_DEBUG_VALUE = "1"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetEnvOrDefault(key string, defaultValue string) string {
	val := os.Getenv(key)

	if val == "" {
		return defaultValue
	}

	return val
}

func GetEnvOrDefaultNum(key string, defaultValue int) int {
	val := os.Getenv(key)

	if val == "" {
		return defaultValue
	}

	i, _ := strconv.ParseInt(val, 10, 32)

	return int(i)
}

func GetEnvOrDefaultFloat(key string, defaultValue float64) float64 {
	val := os.Getenv(key)

	if val == "" {
		return defaultValue
	}

	i, _ := strconv.ParseFloat(val, 64)

	return i
}

func IsDev() bool {
	return os.Getenv(IS_DEV_KEY) == IS_DEV_VALUE
}

func GetGoogleProjectId() string {
	return GetEnvOrDefault(GOOGLE_PROJECT_ID_KEY, "maosaobra-5b2f9")
}

func GetGoogleCredentialsFilepath() string {
	return GetEnvOrDefault(GOOGLE_APPLICATION_CREDENTIALS_KEY, "")
}
