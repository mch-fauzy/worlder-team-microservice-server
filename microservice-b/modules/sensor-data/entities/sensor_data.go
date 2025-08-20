package entities

import (
	"time"

	"gorm.io/gorm"
)

// SensorData represents the sensor data entity
type SensorData struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	SensorValue float64        `json:"sensor_value" gorm:"type:decimal(10,4);not null"`
	SensorType  string         `json:"sensor_type" gorm:"type:varchar(50);not null;index"`
	ID1         string         `json:"id1" gorm:"type:varchar(50);not null;index:idx_id_combination"`
	ID2         int32          `json:"id2" gorm:"not null;index:idx_id_combination"`
	Timestamp   time.Time      `json:"timestamp" gorm:"type:timestamp;not null;index"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName sets the table name for GORM
func (SensorData) TableName() string {
	return "sensor_data"
}
