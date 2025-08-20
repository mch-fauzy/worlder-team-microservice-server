package services

import (
	"context"
	"math"
	"time"

	"github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/dtos"
	"github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/entities"
	"github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/interfaces"
	"github.com/worlder-team/microservice-server/shared/constants"
)

type sensorService struct {
	sensorRepo interfaces.SensorRepositoryInterface
}

// NewSensorService creates a new sensor service
func NewSensorService(sensorRepo interfaces.SensorRepositoryInterface) interfaces.SensorServiceInterface {
	return &sensorService{
		sensorRepo: sensorRepo,
	}
}

func (s *sensorService) CreateSensorData(ctx context.Context, data *entities.SensorData) error {
	return s.sensorRepo.Create(ctx, data)
}

func (s *sensorService) CreateSensorDataBatch(ctx context.Context, data []*entities.SensorData) error {
	return s.sensorRepo.CreateBatch(ctx, data)
}

func (s *sensorService) GetSensorData(ctx context.Context, id uint) (*entities.SensorData, error) {
	return s.sensorRepo.GetByID(ctx, id)
}

func (s *sensorService) GetSensorDataByIDCombination(ctx context.Context, id1 string, id2 int32) ([]*entities.SensorData, error) {
	return s.sensorRepo.GetByIDCombination(ctx, id1, id2)
}

func (s *sensorService) GetSensorDataByDuration(ctx context.Context, from, to time.Time) ([]*entities.SensorData, error) {
	return s.sensorRepo.GetByDuration(ctx, from, to)
}

func (s *sensorService) ListSensorData(ctx context.Context, filter *dtos.SensorDataFilter, pagination *dtos.PaginationParams) (*dtos.PaginatedResponse, error) {
	// Set default pagination if not provided
	if pagination == nil {
		pagination = &dtos.PaginationParams{
			Page:     1,
			PageSize: constants.DefaultPageSize,
			Sort:     "timestamp",
			Order:    "desc",
		}
	}

	// Validate and set defaults
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	if pagination.PageSize < 1 || pagination.PageSize > constants.MaxPageSize {
		pagination.PageSize = constants.DefaultPageSize
	}

	data, total, err := s.sensorRepo.List(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &dtos.PaginatedResponse{
		Data:       data,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *sensorService) UpdateSensorData(ctx context.Context, id uint, data *entities.SensorData) error {
	// Check if record exists
	_, err := s.sensorRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.sensorRepo.Update(ctx, id, data)
}

func (s *sensorService) DeleteSensorData(ctx context.Context, filter *dtos.SensorDataFilter) (int64, error) {
	return s.sensorRepo.Delete(ctx, filter)
}

func (s *sensorService) DeleteSensorDataByID(ctx context.Context, id uint) error {
	// Check if record exists
	_, err := s.sensorRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.sensorRepo.DeleteByID(ctx, id)
}
