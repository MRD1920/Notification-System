package config

import "github.com/joho/godotenv"

func LoadConfig() error {
	err := godotenv.Load()
	return err
}
