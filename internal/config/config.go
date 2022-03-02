package config

import (
	"os"
	"strconv"
)

var serviceVersion = "local"

const (
	redisServerKey   = "REDIS_SERVER"
	redisPasswordKey = "REDIS_PASSWORD"
	port             = "HTTP_PORT"
	jsonLogs         = "JSON_LOGS"
)

type Config struct {
	RedisServer   string
	RedisPassword string
	Port          string
	JSONLogs      bool
}

func New() Config {
	return Config{
		RedisServer:   GetEnvString(redisServerKey, ""),
		RedisPassword: GetEnvString(redisPasswordKey, ""),
		Port:          GetEnvString(port, "8080"),
		JSONLogs:      GetEnvBool(jsonLogs, false),
	}
}

func GetVersion() string {
	return serviceVersion
}

func GetEnvString(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultValue
}

func GetEnvBool(key string, defaultValue bool) bool {
	if val := os.Getenv(key); val != "" {
		bVal, err := strconv.ParseBool(val)
		if err != nil {
			return defaultValue
		}
		return bVal
	}

	return defaultValue
}
