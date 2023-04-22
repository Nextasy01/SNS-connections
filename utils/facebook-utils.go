package utils

import (
	"github.com/joho/godotenv"
)

type FacebookEnvReader struct{}

func (f *FacebookEnvReader) ReadFromEnv() (string, string, error) {
	envFile, err := godotenv.Read(".env")

	if err != nil {
		return "", "", err
	}

	return envFile["Facebook_APP_id"], envFile["Facebook_APP_secret"], nil
}
