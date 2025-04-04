package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

var defaultPort = "7540"

type Environment struct {
	TodoPort string `env:"TODO_PORT" env-required:"true"`
}

// LoadEnv загружает конфигурацию из переменных окружения, в случае ошибки или если переменные не найдены то присваивает дефолтные значения
func LoadEnv() Environment {
	var env Environment

	_ = godotenv.Load() // Игнорируем ошибку если .env не найден

	if err := cleanenv.ReadEnv(&env); err != nil {
		env.TodoPort = defaultPort
	}

	if env.TodoPort == "" {
		env.TodoPort = defaultPort
	}

	return env
}
