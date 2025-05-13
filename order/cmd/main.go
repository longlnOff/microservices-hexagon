package main

import (
	"fmt"
	"go.uber.org/zap"
)

func main() {
	logger := logger.CreateLogger()
	defer logger.Sync()

	cfg, err := configuration.LoadConfig(".")
	if err != nil {
		logger.Fatal(err.Error())
	}

	database, err := postgres.New(
		fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.Database.ENGINE,
			cfg.Database.USER,
			cfg.Database.PASSWORD,
			cfg.Database.HOST,
			cfg.Database.PORT,
			cfg.Database.DB_NAME,
			cfg.Database.DB_SSL_MODE),
		cfg.Database.DB_MAX_OPEN_CONNS,
		cfg.Database.DB_MAX_IDLE_CONNS,
		cfg.Database.DB_MAX_IDLE_TIME,
	)
	if err != nil {
		logger.Panic(err.Error())
	}
	defer database.Close()
	logger.Info("Connected to database.", zap.String("url", fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Database.USER, cfg.Database.PASSWORD, cfg.Database.HOST, cfg.Database.PORT, cfg.Database.DB_NAME)))
}
