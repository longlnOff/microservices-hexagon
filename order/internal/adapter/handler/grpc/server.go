package grpc

import (
	"fmt"
	"net"

	"github.com/longlnOff/microservices-hexagon/order/configuration"
	"github.com/longlnOff/microservices-hexagon/order/internal/port"
	"github.com/longlnOff/microservices-proto/golang/order"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	configuration configuration.Configuration
	logger *zap.Logger
	service port.OrderService
	port 	int
	server *grpc.Server
	order.UnimplementedOrderServer
}

func NewGRPCServer(config configuration.Configuration, logger *zap.Logger, service port.OrderService, port int) *GRPCServer {
	return &GRPCServer{
		configuration: config,
		logger:        logger,
		service:       service,
		port:          port,
	}
}

func (g GRPCServer) Run() {
	var err error

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		g.logger.Fatal("failed to listen on port", zap.Int("port", g.port), zap.Error(err))
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	g.server = grpcServer
	order.RegisterOrderServer(grpcServer, g)
	if g.configuration.GRPCServer.GRPC_ENVIRONMENT == "development" {
		reflection.Register(grpcServer)
	}

	g.logger.Info("starting order service on port", zap.Int("port", g.port))
	if err := grpcServer.Serve(listen); err != nil {
		g.logger.Fatal("failed to serve grpc on port", zap.Int("port", g.port), zap.Error(err))
	}
}

func (g GRPCServer) Stop() {
	g.server.Stop()
}
