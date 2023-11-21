package banner

import (
	"context"
	"fmt"
	"github.com/dianapovarnitsina/banners-rotation/internal/server/pb"
	"net"
	"os"

	"github.com/dianapovarnitsina/banners-rotation/interfaces"
	"github.com/dianapovarnitsina/banners-rotation/internal/config"
	"github.com/dianapovarnitsina/banners-rotation/internal/logger"
	internalgrpc "github.com/dianapovarnitsina/banners-rotation/internal/server/grpc"
	"github.com/dianapovarnitsina/banners-rotation/internal/storage/sql"
	"google.golang.org/grpc"
)

type App struct {
	logger       interfaces.Logger
	storage      interfaces.Storage
	serverGRPC   *grpc.Server
	grpcShutdown chan struct{} // Канал для сигнала завершения gRPC сервера
	//queue       Queue
	//eventsQueue string
}

func NewApp(ctx context.Context, conf *config.BannerConfig) (*App, error) {
	app := &App{}

	logger := logger.New(conf.Logger.Level, os.Stdout)
	app.logger = logger

	//Инициализация хранилища данных.
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

	//Инициализация gRPC-сервера.
	app.serverGRPC = grpc.NewServer(
		grpc.UnaryInterceptor(internalgrpc.NewLoggingInterceptor(logger).UnaryServerInterceptor),
	)

	api := internalgrpc.NewEventServiceServer(app.storage)
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
		close(app.grpcShutdown) // Отправляем сигнал о завершении работы gRPC сервера.
	}()

	return app, nil
}

func (a *App) GetGrpcServerShutdownSignal() <-chan struct{} {
	return a.grpcShutdown
}
