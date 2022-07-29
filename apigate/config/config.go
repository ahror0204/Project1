package config

import (
	"os"

	"github.com/spf13/cast"
)

type Config struct {
	Environment string // Develop, staging, production

	UserServiceHost string
	UserServicePort int

	PostgresHost string
	PostgresPort int
	PostgresUser string
	PostgresPassword string
	PostgresDatabase string

	//context timeout in seconds
	CtxTimeout int
	RedisHost  string
	RedisPort  int

	LogeLevel string
	HTTPPort  string

	SigningKey string
}

func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))

	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "1"))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DB", "template"))

	c.LogeLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":8080"))

	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "127.0.0.1"))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 9000))

	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "localhost"))
	c.RedisPort = cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379))

	c.SigningKey = cast.ToString(getOrReturnDefault("SIGNING_KEY", "bzqymwhbwgholtdyzjvqaycuxwnmeqczzosvmafrjfskmepquudmdktutkyzowntnvwurvkxywkpxsexhkkwcqgsbbbxlqyuklcrbypczsfhwejwqebsxqmprueopdexwdmukhfkujxhjeecfiwwspjgbgcgowew"))
	
	c.CtxTimeout = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
