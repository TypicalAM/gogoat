package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Site  string `json:"site"`
	Token string `json:"token"`
}

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("loading .env file: %w", err)
	}

	sitePrefix, ok := os.LookupEnv("SITE_PREFIX")
	if !ok {
		return nil, fmt.Errorf("SITE not found in .env")
	}

	token, ok := os.LookupEnv("TOKEN")
	if !ok {
		return nil, fmt.Errorf("TOKEN not found in .env")
	}

	return &Config{
		Site:  fmt.Sprintf("https://%s.goatcounter.com", sitePrefix),
		Token: token,
	}, nil
}
