package dtos

import "time"

// SensorDataFilter represents filter criteria for sensor data queries
type SensorDataFilter struct {
	SensorType *string    `json:"sensor_type,omitempty"`
	ID1        *string    `json:"id1,omitempty"`
	ID2        *int32     `json:"id2,omitempty"`
	FromTime   *time.Time `json:"from_time,omitempty"`
	ToTime     *time.Time `json:"to_time,omitempty"`
	MinValue   *float64   `json:"min_value,omitempty"`
	MaxValue   *float64   `json:"max_value,omitempty"`
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	Sort     string `json:"sort" query:"sort"`
	Order    string `json:"order" query:"order"`
}

// PaginatedResponse represents paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}
