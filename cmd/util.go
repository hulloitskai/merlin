package cmd

import (
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/joho/godotenv"
	ess "github.com/unixpickle/essentials"
)

func loadEnv() {
	var err error
	if os.Getenv("GO_ENV") == "development" {
		err = godotenv.Load(".env", ".env.local")
	} else {
		err = godotenv.Load("/secrets/env/.env")
	}

	if (err != nil) &&
		!strings.Contains(err.Error(), "no such file or directory") {
		ess.Die("Error while reading .env file:", err)
	}
}

func buildLogger() (*zap.SugaredLogger, error) {
	var (
		raw *zap.Logger
		err error
	)
	if os.Getenv("GO_ENV") == "development" {
		raw, err = zap.NewDevelopment()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.Encoding = "console"
		raw, err = cfg.Build()
	}
	if err != nil {
		return nil, err
	}
	return raw.Sugar(), nil
}
