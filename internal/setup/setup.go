package setup

import "os"

type EnvData struct {
	SECRET            string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
}

func SetupEnv() *EnvData {
	return &EnvData{
		SECRET: os.Getenv("SECRET"),
		POSTGRES_USER: os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_DB: os.Getenv("POSTGRES_DB"),
	}
}