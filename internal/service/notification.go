package service

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"go.uber.org/zap"
)

type NotificationService struct {
	log        *zap.Logger
	repository repository.INotificationRepository
}

type INotificationService interface {
	CreateNotification(request dto.CreateNotificationRequest) (dto.NotificationResponse, error)
	GetNotifications(filter dto.GetNotificationFilter) ([]dto.NotificationResponse, error)
	MarkNotifications(request dto.MarkNotificationRequest, userId uuid.UUID) ([]dto.NotificationResponse, error)
}

func NewNotificationService(log *zap.Logger, repository repository.INotificationRepository) INotificationService {
	return &NotificationService{
		log:        log,
		repository: repository,
	}
}

func (s *NotificationService) CreateNotification(request dto.CreateNotificationRequest) (dto.NotificationResponse, error) {
	s.repository.UseTx(false)

	data := entity.Notification{
		Description: request.Description,
	}

	if request.UserId != nil {
		data.UserId = uuid.NullUUID{UUID: uuid.MustParse(*request.UserId), Valid: true}
	}

	if request.CageId != nil {
		data.CageId = sql.NullInt64{Int64: int64(*request.CageId), Valid: true}
	}

	if request.WarehouseId != nil {
		data.WarehouseId = sql.NullInt64{Int64: int64(*request.WarehouseId), Valid: true}
	}

	if request.StoreId != nil {
		data.StoreId = sql.NullInt64{Int64: int64(*request.StoreId), Valid: true}
	}

	err := s.repository.CreateNotification(&data)
	if err != nil {
		s.log.Error("failed create notification", zap.Error(err))
		return dto.NotificationResponse{}, err
	}

	return dto.NotificationResponse{
		Id:          data.Id,
		Description: data.Description,
		IsMarked:    data.IsMarked,
	}, nil
}

func (s *NotificationService) GetNotifications(filter dto.GetNotificationFilter) ([]dto.NotificationResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetNotifications(filter)
	if err != nil {
		s.log.Error("failed get notifications", zap.Error(err))
		return nil, err
	}

	response := make([]dto.NotificationResponse, 0)
	for _, d := range data {
		response = append(response, dto.NotificationResponse{
			Id:                   d.Id,
			Description:          d.Description,
			IsMarked:             d.IsMarked,
			NotificationContexts: d.NotificationContexts,
		})
	}

	return response, nil
}

func (s *NotificationService) MarkNotifications(request dto.MarkNotificationRequest, userId uuid.UUID) ([]dto.NotificationResponse, error) {
	s.repository.UseTx(false)

	err := s.repository.UpdateNotificationAsMarked(request.Ids)
	if err != nil {
		s.log.Error("failed update notification as marked", zap.Error(err))
		return nil, err
	}

	data, err := s.repository.GetNotificationByIds(request.Ids)
	if err != nil {
		s.log.Error("failed get notifications", zap.Error(err))
		return nil, err
	}

	response := make([]dto.NotificationResponse, 0)
	for _, d := range data {
		response = append(response, dto.NotificationResponse{
			Id:          d.Id,
			Description: d.Description,
			IsMarked:    d.IsMarked,
		})
	}

	return response, nil
}
