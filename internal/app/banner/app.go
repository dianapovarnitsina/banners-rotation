package banner

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/dianapovarnitsina/banners-rotation/interfaces"
	"github.com/dianapovarnitsina/banners-rotation/internal/config"
	"github.com/dianapovarnitsina/banners-rotation/internal/logger"
	"github.com/dianapovarnitsina/banners-rotation/internal/rmq"
	internalgrpc "github.com/dianapovarnitsina/banners-rotation/internal/server/grpc"
	"github.com/dianapovarnitsina/banners-rotation/internal/server/pb"
	"github.com/dianapovarnitsina/banners-rotation/internal/storage/sql"
	"google.golang.org/grpc"
)

type App struct {
	logger     interfaces.Logger
	storage    interfaces.Storage
	serverGRPC *grpc.Server
}

func NewApp(ctx context.Context, conf *config.BannerConfig) (*App, error) {
	app := &App{}

	logger := logger.New(conf.Logger.Level, os.Stdout)
	app.logger = logger

	// Инициализация хранилища данных.
	psqlStorage := new(sql.Storage)
	if err := psqlStorage.Connect(
		ctx,
		conf.Database.Port,
		conf.Database.Host,
		conf.Database.Username,
		conf.Database.Password,
		conf.Database.Dbname,
	); err != nil {
		return nil, fmt.Errorf("cannot connect to PostgreSQL: %w", err)
	}
	err := psqlStorage.Migrate(ctx, conf.Storage.Migration)
	if err != nil {
		return nil, fmt.Errorf("migration did not work out: %w", err)
	}
	app.storage = psqlStorage

	// Инициализация RMQ.
	URI := fmt.Sprintf("%s://%s:%s@%s:%d/", //"amqp://guest:guest@localhost:5672/"
		conf.RMQ.RABBITMQ_PROTOCOL,
		conf.RMQ.RABBITMQ_USERNAME,
		conf.RMQ.RABBITMQ_PASSWORD,
		conf.RMQ.RABBITMQ_HOST,
		conf.RMQ.RABBITMQ_PORT,
	)
	logger.Info("URI: ", URI)

	eventsProdMq, err := rmq.New(
		URI,
		conf.Queues.Events.ExchangeName,
		conf.Queues.Events.ExchangeType,
		conf.Queues.Events.QueueName,
		conf.Queues.Events.BindingKey,
		conf.RMQ.ReConnect.MaxElapsedTime,
		conf.RMQ.ReConnect.InitialInterval,
		conf.RMQ.ReConnect.Multiplier,
		conf.RMQ.ReConnect.MaxInterval,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize RMQ for scheduler: %w", err)
	}

	if err := eventsProdMq.Init(ctx); err != nil {
		logger.Error("RMQ initialization failed: %v", err)
	}

	// Инициализация gRPC-сервера.
	app.serverGRPC = grpc.NewServer(
		grpc.UnaryInterceptor(internalgrpc.NewLoggingInterceptor(logger).UnaryServerInterceptor),
	)

	api := internalgrpc.NewEventServiceServer(app.storage, eventsProdMq, logger)
	pb.RegisterBannerServiceServer(app.serverGRPC, api)

	grpcListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.GRPC.Host, conf.GRPC.Port))
	if err != nil {
		logger.Error("Failed to listen: %v", err)
	}

	go func() {
		logger.Info("Starting gRPC server on port %s", fmt.Sprintf("%s:%d", conf.GRPC.Host, conf.GRPC.Port))
		if err := app.serverGRPC.Serve(grpcListener); err != nil {
			logger.Error("gRPC server failed: %v", err)
		}
	}()

	// Ожидание сигнала завершения работы сервера.
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		// Получен сигнал завершения - остановка gRPC-сервера.
		app.serverGRPC.GracefulStop()
		logger.Info("gRPC server stopped")
	}()

	return app, nil
}
