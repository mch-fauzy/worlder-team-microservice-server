package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/dtos"
	"github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/entities"
	"github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/interfaces"
	"gorm.io/gorm"
)

type sensorRepository struct {
	db *gorm.DB
}

// NewSensorRepository creates a new sensor repository
func NewSensorRepository(db *gorm.DB) interfaces.SensorRepositoryInterface {
	return &sensorRepository{
		db: db,
	}
}

func (r *sensorRepository) Create(ctx context.Context, data *entities.SensorData) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *sensorRepository) CreateBatch(ctx context.Context, data []*entities.SensorData) error {
	return r.db.WithContext(ctx).CreateInBatches(data, 100).Error
}

func (r *sensorRepository) GetByID(ctx context.Context, id uint) (*entities.SensorData, error) {
	var data entities.SensorData
	err := r.db.WithContext(ctx).First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *sensorRepository) GetByIDCombination(ctx context.Context, id1 string, id2 int32) ([]*entities.SensorData, error) {
	var data []*entities.SensorData
	err := r.db.WithContext(ctx).Where("id1 = ? AND id2 = ?", id1, id2).Find(&data).Error
	return data, err
}

func (r *sensorRepository) GetByDuration(ctx context.Context, from, to time.Time) ([]*entities.SensorData, error) {
	var data []*entities.SensorData
	err := r.db.WithContext(ctx).Where("timestamp BETWEEN ? AND ?", from, to).Find(&data).Error
	return data, err
}

func (r *sensorRepository) List(ctx context.Context, filter *dtos.SensorDataFilter, pagination *dtos.PaginationParams) ([]*entities.SensorData, int64, error) {
	var data []*entities.SensorData
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.SensorData{})

	// Apply filters
	if filter != nil {
		if filter.SensorType != nil {
			query = query.Where("sensor_type = ?", *filter.SensorType)
		}
		if filter.ID1 != nil {
			query = query.Where("id1 = ?", *filter.ID1)
		}
		if filter.ID2 != nil {
			query = query.Where("id2 = ?", *filter.ID2)
		}
		if filter.FromTime != nil {
			query = query.Where("timestamp >= ?", *filter.FromTime)
		}
		if filter.ToTime != nil {
			query = query.Where("timestamp <= ?", *filter.ToTime)
		}
		if filter.MinValue != nil {
			query = query.Where("sensor_value >= ?", *filter.MinValue)
		}
		if filter.MaxValue != nil {
			query = query.Where("sensor_value <= ?", *filter.MaxValue)
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.PageSize
		query = query.Offset(offset).Limit(pagination.PageSize)

		// Apply sorting
		if pagination.Sort != "" {
			order := pagination.Order
			if order != "asc" && order != "desc" {
				order = "asc"
			}
			query = query.Order(fmt.Sprintf("%s %s", pagination.Sort, order))
		} else {
			query = query.Order("created_at desc")
		}
	}

	err := query.Find(&data).Error
	return data, total, err
}

func (r *sensorRepository) Update(ctx context.Context, id uint, data *entities.SensorData) error {
	// Use Select to explicitly update all updatable fields, including zero values
	return r.db.WithContext(ctx).Model(&entities.SensorData{}).Where("id = ?", id).
		Select("sensor_value", "sensor_type", "timestamp").Updates(data).Error
}

func (r *sensorRepository) Delete(ctx context.Context, filter *dtos.SensorDataFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&entities.SensorData{})

	// Apply filters
	if filter != nil {
		if filter.SensorType != nil {
			query = query.Where("sensor_type = ?", *filter.SensorType)
		}
		if filter.ID1 != nil {
			query = query.Where("id1 = ?", *filter.ID1)
		}
		if filter.ID2 != nil {
			query = query.Where("id2 = ?", *filter.ID2)
		}
		if filter.FromTime != nil {
			query = query.Where("timestamp >= ?", *filter.FromTime)
		}
		if filter.ToTime != nil {
			query = query.Where("timestamp <= ?", *filter.ToTime)
		}
		if filter.MinValue != nil {
			query = query.Where("sensor_value >= ?", *filter.MinValue)
		}
		if filter.MaxValue != nil {
			query = query.Where("sensor_value <= ?", *filter.MaxValue)
		}
	}

	result := query.Delete(&entities.SensorData{})
	return result.RowsAffected, result.Error
}

func (r *sensorRepository) DeleteByID(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.SensorData{}, id).Error
}
