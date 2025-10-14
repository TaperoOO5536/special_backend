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
	"google.golang.org/protobuf/encoding/protojson"
)

type Config struct {
	GrpcPort     string
	HttpPort     string
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
	eventRepo := repository.NewEventRepository(db)
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	userEventRepo := repository.NewUserEventRepository(db)
	
	itemService := service.NewItemService(itemRepo)
	eventService := service.NewEventService(eventRepo)
	userServive := service.NewUserService(userRepo, config.GetToken())
	orderService := service.NewOrderService(orderRepo, config.GetToken())
	userEventService := service.NewUserEventService(userEventRepo, eventRepo, config.GetToken())

	itemServiceHandler := api.NewItemServiceHandler(itemService)
	eventServiceHandler := api.NewEventServiceHandler(eventService)
	userServiceHandler := api.NewUserServiceHandler(userServive)
	orderServiceHandler := api.NewOrderServiceHandler(orderService)
	userEventServiceHandler := api.NewUserEventServiceHandler(userEventService)

	handler := api.NewHandler(itemServiceHandler, eventServiceHandler, userServiceHandler, orderServiceHandler, userEventServiceHandler)

	grpcServer := grpc.NewServer()

	pb.RegisterSpecialAppServiceServer(grpcServer, handler)

	reflection.Register(grpcServer)

	l, err := net.Listen("tcp", ":"+a.config.GrpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	ctx = context.Background()
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
            MarshalOptions: protojson.MarshalOptions{
                EmitUnpopulated: true,
            },
            UnmarshalOptions: protojson.UnmarshalOptions{
                DiscardUnknown: true,
            },
        }),
	)
	err = pb.RegisterSpecialAppServiceHandlerFromEndpoint(
		ctx,
		gwmux,
		"localhost:"+a.config.GrpcPort,
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
		Addr:    ":"+a.config.HttpPort,
		Handler: corshandler,
	}

	serverError := make(chan error, 1)
	go func() {
		log.Printf("Starting gRPC server on port %s", a.config.GrpcPort)
		serverError <- grpcServer.Serve(l)
	}()

	go func() {
		log.Printf("Starting http server on port %s", a.config.HttpPort)
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
		log.Println("shutting down http server...")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown error: %v", err)
		}
		log.Println("http server stopped")
		log.Println("gRPC server stopped")
		return nil
	}

}