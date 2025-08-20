package entities

import "time"

// SensorData represents sensor data structure
type SensorData struct {
	SensorValue float64   `json:"sensor_value"`
	SensorType  string    `json:"sensor_type"`
	ID1         string    `json:"id1"`
	ID2         int32     `json:"id2"`
	Timestamp   time.Time `json:"timestamp"`
}
