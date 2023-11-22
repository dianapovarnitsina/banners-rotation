package internalgrpc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dianapovarnitsina/banners-rotation/interfaces"
	"github.com/dianapovarnitsina/banners-rotation/internal/logger"
	"github.com/dianapovarnitsina/banners-rotation/internal/rmq"
	"github.com/dianapovarnitsina/banners-rotation/internal/server/pb"
	"github.com/dianapovarnitsina/banners-rotation/internal/storage"
	"github.com/streadway/amqp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServiceServer struct {
	storage      interfaces.Storage
	eventsProdMq *rmq.Rmq
	logger       *logger.Logger
	pb.UnimplementedBannerServiceServer
}

func NewEventServiceServer(storage interfaces.Storage, eventsProdMq *rmq.Rmq, log *logger.Logger) *ServiceServer {
	return &ServiceServer{
		storage:      storage,
		eventsProdMq: eventsProdMq,
		logger:       log,
	}
}

func (s *ServiceServer) AddBanner(ctx context.Context, req *pb.AddBannerRequest) (*pb.AddBannerResponse, error) {
	bannerID := int(req.GetBannerId())
	slotID := int(req.GetSlotId())

	// Проверка на несуществующий баннер
	if !s.bannerExists(ctx, bannerID) {
		return nil, status.Errorf(codes.NotFound, "specified banner does not exist")
	}

	// Проверка на несуществующий слот
	if !s.slotExists(ctx, slotID) {
		return nil, status.Errorf(codes.NotFound, "specified slot does not exist")
	}

	// Проверка на повторное добавление баннера в слот
	if exists, err := s.checkDuplicateBannerSlot(ctx, bannerID, slotID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check duplicate: %v", err)
	} else if exists {
		return nil, status.Errorf(codes.AlreadyExists, "banner is already assigned to the slot")
	}

	//добавление записи
	if err := s.storage.AddBanner(ctx, bannerID, slotID); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add banner: %v", err)
	}

	return &pb.AddBannerResponse{Message: "Banner added successfully"}, nil
}

func (s *ServiceServer) checkDuplicateBannerSlot(ctx context.Context, bannerID, slotID int) (bool, error) {
	return s.storage.IsBannerAssignedToSlot(ctx, bannerID, slotID)
}

func (s *ServiceServer) bannerExists(ctx context.Context, bannerID int) bool {
	return s.storage.BannerExists(ctx, bannerID)
}

func (s *ServiceServer) slotExists(ctx context.Context, slotID int) bool {
	return s.storage.SlotExists(ctx, slotID)
}

func (s *ServiceServer) userGroupExists(ctx context.Context, userGroupId int) bool {
	return s.storage.UserGroupExists(ctx, userGroupId)
}

func (s *ServiceServer) RemoveBanner(ctx context.Context, req *pb.RemoveBannerRequest) (*pb.RemoveBannerResponse, error) {
	if err := s.storage.RemoveBanner(ctx, int(req.GetBannerId()), int(req.GetSlotId())); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to remove banner: %v", err)
	}
	return &pb.RemoveBannerResponse{Message: "Banner removed successfully"}, nil
}

func (s *ServiceServer) ClickBanner(ctx context.Context, req *pb.ClickBannerRequest) (*pb.ClickBannerResponse, error) {
	bannerID := int(req.GetBannerId())
	slotID := int(req.GetSlotId())
	userGroupId := int(req.GetUsergroupId())

	// Проверка на несуществующий баннер
	if !s.bannerExists(ctx, bannerID) {
		return nil, status.Errorf(codes.NotFound, "specified banner does not exist")
	}

	// Проверка на несуществующий слот
	if !s.slotExists(ctx, slotID) {
		return nil, status.Errorf(codes.NotFound, "specified slot does not exist")
	}

	// Проверка на несуществующий группу
	if !s.userGroupExists(ctx, userGroupId) {
		return nil, status.Errorf(codes.NotFound, "specified userGroup does not exist")
	}

	click, err := s.storage.ClickBanner(ctx, bannerID, slotID, userGroupId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to click banner: %v", err)
	}

	// Отправка уведомления в очередь
	notification := createClickNotification(click)
	err = s.sendNotification(notification)

	if err != nil {
		s.logger.Error("Failed to RMQ: %v", err)
	}
	return &pb.ClickBannerResponse{Message: "Banner clicked successfully"}, nil
}

func (s *ServiceServer) PickBanner(ctx context.Context, req *pb.PickBannerRequest) (*pb.PickBannerResponse, error) {
	slotID := int(req.GetSlotId())
	userGroupId := int(req.GetUsergroupId())

	impress, bannerID, err := s.storage.PickBanner(ctx, slotID, userGroupId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to pick banner: %v", err)
	}

	// Отправка уведомления в очередь
	notification := createImpressNotification(impress)
	err = s.sendNotification(notification)

	if err != nil {
		s.logger.Error("Failed to RMQ: %v", err)
	}

	return &pb.PickBannerResponse{BannerId: int32(bannerID), Message: "Banner picked successfully"}, nil
}

func (s *ServiceServer) sendNotification(notification storage.Notification) error {
	notificationJSON, err := serializeNotification(notification)
	if err != nil {
		s.logger.Error("Failed to serialize notification: %v", err)
		return err
	}

	err = publishNotificationToRMQ(notificationJSON, s.eventsProdMq)
	if err != nil {
		s.logger.Error("Failed to publish notification to RMQ: %v", err)
		return err
	}

	s.logger.Info("Sent a notification to queue RabbitMQ: %s %s", string(notificationJSON), time.Now().Format("2006-01-02 15:04"))
	return nil
}

func createClickNotification(сlick *storage.Click) storage.Notification {
	notification := storage.Notification{
		TypeEvent:   "click",
		SlotId:      сlick.SlotID,
		BannerId:    сlick.BannerID,
		UsergroupId: сlick.UserGroupID,
		DateTime:    сlick.CreatedAt,
	}
	return notification
}

func createImpressNotification(сlick *storage.Impress) storage.Notification {
	notification := storage.Notification{
		TypeEvent:   "impress",
		SlotId:      сlick.SlotID,
		BannerId:    сlick.BannerID,
		UsergroupId: сlick.UserGroupID,
		DateTime:    сlick.CreatedAt,
	}
	return notification
}

func serializeNotification(notification storage.Notification) ([]byte, error) {
	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		return nil, err
	}
	return notificationJSON, nil
}

func publishNotificationToRMQ(notificationJSON []byte, rmq *rmq.Rmq) error {
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        notificationJSON,
	}

	return rmq.Publish(msg)
}
