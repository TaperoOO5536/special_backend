package app

import (
	"context"
	"fmt"
	"log"

	"net"
	"net/http"
	"time"

	"github.com/TaperoOO5536/special_backend/internal/config"
	"github.com/TaperoOO5536/special_backend/internal/repository"
	"github.com/TaperoOO5536/special_backend/internal/service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"

	"os"
	"os/signal"
	"syscall"

	"github.com/TaperoOO5536/special_backend/internal/api"
	pb "github.com/TaperoOO5536/special_backend/pkg/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	Port string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Dsn          string
}

type App struct {
	config *Config
}

func New(cfg *Config) *App {
	return &App {
		config: cfg,
	}
}

func (a *App) Start(ctx context.Context) error {
	db := config.NewDBClient(a.config.Dsn)

	itemRepo := repository.NewItemRepository(db)
	iventRepo := repository.NewIventRepository(db)
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	userIventRepo := repository.NewUserIventRepository(db)
	
	itemService := service.NewItemService(itemRepo)
	iventService := service.NewIventService(iventRepo)
	userServive := service.NewUserService(userRepo, config.GetToken())
	orderService := service.NewOrderService(orderRepo, config.GetToken())
	userIventService := service.NewUserIventService(userIventRepo, iventRepo, config.GetToken())

	itemServiceHandler := api.NewItemServiceHandler(itemService)
	iventServiceHandler := api.NewIventServiceHandler(iventService)
	userServiceHandler := api.NewUserServiceHandler(userServive)
	orderServiceHandler := api.NewOrderServiceHandler(orderService)
	userIventServiceHandler := api.NewUserIventServiceHandler(userIventService)

	handler := api.NewHandler(itemServiceHandler, iventServiceHandler, userServiceHandler, orderServiceHandler, userIventServiceHandler)

	grpcServer := grpc.NewServer()

	pb.RegisterSpecialAppServiceServer(grpcServer, handler)

	reflection.Register(grpcServer)

	l, err := net.Listen("tcp", ":"+a.config.Port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	ctx = context.Background()
	gwmux := runtime.NewServeMux()
	err = pb.RegisterSpecialAppServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		"localhost:"+a.config.Port,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		return fmt.Errorf("failed to register gateway: %v", err)
	}

	c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://192.168.1.212", "http://localhost:5173"},
        AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept", "X-Requested-With"},
        AllowCredentials: true,
        MaxAge:           300,
				Debug:            true,
    })

	corshandler := c.Handler(gwmux)

	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: corshandler,
	}

	serverError := make(chan error, 1)
	go func() {
		log.Printf("Starting gRPC server on port%s", a.config.Port)
		serverError <- grpcServer.Serve(l)
	}()

	go func() {
		log.Printf("Starting http server on port :8081")
		serverError <- httpServer.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <- serverError:
		return fmt.Errorf("grpc server error %v", err)
	case <-shutdown:
		log.Println("shutting down gRPC server...")
		grpcServer.GracefulStop()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown error: %v", err)
		}
		log.Println("gRPC server stopped")
		return nil
	}

}