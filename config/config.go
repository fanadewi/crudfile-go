package config

import (
	"os"
	"strconv"
)

type Config struct {
	DebugMode   bool
	Port        int
	CloudName   string
	CloudKey    string
	CloudSecret string
}

func New() *Config {
	return &Config{
		DebugMode:   getEnvAsBool("DEBUG_MODE", true),
		Port:        getEnvAsInt("PORT", 8080),
		CloudName:   getEnv("CLOUD_NAME", ""),
		CloudKey:    getEnv("CLOUD_KEY", ""),
		CloudSecret: getEnv("CLOUD_SECRET", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}
