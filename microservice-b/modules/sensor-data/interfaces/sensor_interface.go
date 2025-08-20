package interfaces

import (
	"context"
	"time"

	"github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/dtos"
	"github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/entities"
)

// SensorRepositoryInterface defines the interface for sensor data repository
type SensorRepositoryInterface interface {
	Create(ctx context.Context, data *entities.SensorData) error
	CreateBatch(ctx context.Context, data []*entities.SensorData) error
	GetByID(ctx context.Context, id uint) (*entities.SensorData, error)
	GetByIDCombination(ctx context.Context, id1 string, id2 int32) ([]*entities.SensorData, error)
	GetByDuration(ctx context.Context, from, to time.Time) ([]*entities.SensorData, error)
	List(ctx context.Context, filter *dtos.SensorDataFilter, pagination *dtos.PaginationParams) ([]*entities.SensorData, int64, error)
	Update(ctx context.Context, id uint, data *entities.SensorData) error
	Delete(ctx context.Context, filter *dtos.SensorDataFilter) (int64, error)
	DeleteByID(ctx context.Context, id uint) error
}

// SensorServiceInterface defines the interface for sensor service
type SensorServiceInterface interface {
	CreateSensorData(ctx context.Context, data *entities.SensorData) error
	CreateSensorDataBatch(ctx context.Context, data []*entities.SensorData) error
	GetSensorData(ctx context.Context, id uint) (*entities.SensorData, error)
	GetSensorDataByIDCombination(ctx context.Context, id1 string, id2 int32) ([]*entities.SensorData, error)
	GetSensorDataByDuration(ctx context.Context, from, to time.Time) ([]*entities.SensorData, error)
	ListSensorData(ctx context.Context, filter *dtos.SensorDataFilter, pagination *dtos.PaginationParams) (*dtos.PaginatedResponse, error)
	UpdateSensorData(ctx context.Context, id uint, data *entities.SensorData) error
	DeleteSensorData(ctx context.Context, filter *dtos.SensorDataFilter) (int64, error)
	DeleteSensorDataByID(ctx context.Context, id uint) error
}
