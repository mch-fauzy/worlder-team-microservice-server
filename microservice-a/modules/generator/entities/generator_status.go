package entities

import "time"

// GeneratorStatus represents the current status of the generator
type GeneratorStatus struct {
	IsRunning     bool          `json:"is_running"`
	SensorType    string        `json:"sensor_type"`
	Frequency     time.Duration `json:"frequency"`
	LastGenerated time.Time     `json:"last_generated,omitempty"`
	TotalSent     int64         `json:"total_sent"`
	Errors        int64         `json:"errors"`
}
