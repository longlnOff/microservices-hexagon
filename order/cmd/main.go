package main

import (
	"fmt"
	"github.com/longlnOff/microservices-hexagon/order/configuration"
	"github.com/longlnOff/microservices-hexagon/order/internal/adapter/handler/grpc"
	"github.com/longlnOff/microservices-hexagon/order/internal/adapter/logger"
	"github.com/longlnOff/microservices-hexagon/order/internal/adapter/storage/cache"
	"github.com/longlnOff/microservices-hexagon/order/internal/adapter/storage/postgres"
	"github.com/longlnOff/microservices-hexagon/order/internal/adapter/storage/postgres/repository"
	"github.com/longlnOff/microservices-hexagon/order/internal/core/service"
	"go.uber.org/zap"
)

func main() {
	logger := logger.CreateLogger()
	defer logger.Sync()

	useAbsPath := true
	cfg, err := configuration.LoadConfig(".", useAbsPath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// database connection
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
	defer database.DB.Close()
	logger.Info("Connected to database.", zap.String("url", fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.Database.USER, cfg.Database.PASSWORD, cfg.Database.HOST, cfg.Database.PORT, cfg.Database.DB_NAME)))
	// migrate up or not
	if cfg.Database.MIGRATE_UP {
		err = database.MigrateUpTo(cfg.Database.MIGRATE_VERSION)
		if err != nil {
			logger.Fatal(err.Error())
		}
	}

	var cacheClient *cache.Redis
	// cache connection
	if cfg.Cache.CACHE_ENABLED {
		cacheClient = cache.NewValkeyClient(
			cfg.Cache.CACHE_ADDRESS,
			cfg.Cache.CACHE_PASSWORD,
			cfg.Cache.CACHE_DATABASE,
		)
		defer cacheClient.Close()
		if cacheClient == nil {
			logger.Warn("Can't connect to cache.")
		} else {
			logger.Info("Connected to cache at", zap.String("address", cfg.Cache.CACHE_ADDRESS))
		}
	}


	orderRepository := repository.NewOrderRepository(database)
	orderService := service.NewOrderService(orderRepository, cacheClient)
	grpcServer := grpc.NewGRPCServer(
		cfg,
		logger,
		orderService,
		cfg.GRPCServer.GRPC_SERVER_PORT,
	)		
	grpcServer.Run()
}
