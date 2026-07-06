package config

import (
	"errors"
	"os"
)

type Config struct {
	GithubToken string
}

func LoadConfig() (*Config,error){
	token := os.Getenv("GITHUB_TOKEN")
	if token == ""{
		return nil, errors.New("Authentication missing: Please set the GITHUB_TOKEN environment variable")
	}
	return &Config{
		GithubToken: token,
	},nil
}