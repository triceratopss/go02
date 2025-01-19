package config

import (
	"github.com/caarlos0/env/v11"
)

type EnvironmentVariables struct {
	Env        string `env:"ENV" envDefault:"development"`
	DBHost     string `env:"DB_HOST,notEmpty"`
	DBPort     string `env:"DB_PORT,notEmpty"`
	DBName     string `env:"DB_NAME,notEmpty"`
	DBUser     string `env:"DB_USER,notEmpty"`
	DBPassword string `env:"DB_PASSWORD,notEmpty"`
}

var variables EnvironmentVariables

func Init() error {
	var err error
	variables, err = env.ParseAs[EnvironmentVariables]()
	if err != nil {
		return err
	}

	return nil
}

func Get() EnvironmentVariables {
	return variables
}
