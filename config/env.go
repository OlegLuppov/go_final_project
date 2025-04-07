package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Environment struct {
	TodoPort   string `env:"TODO_PORT" env-default:"7540"`
	TodoDbFile string `env:"TODO_DBFILE" env-default:"./scheduler.db"`
}

// LoadEnv загружает конфигурацию из переменных окружения
func LoadEnv() (Environment, error) {
	var env Environment

	godotenv.Load() // Игнорируем ошибку если .env не найден

	if err := cleanenv.ReadEnv(&env); err != nil {
		return env, err
	}

	return env, nil
}
