package dtos

// FrequencyRequest represents frequency change request
type FrequencyRequest struct {
	Frequency string `json:"frequency" validate:"required"`
}

// FrequencyResponse represents frequency response
type FrequencyResponse struct {
	Frequency string `json:"frequency"`
	Duration  string `json:"duration"`
}
