package interfaces

import (
	"context"

	"github.com/worlder-team/microservice-server/microservice-a/modules/generator/entities"
)

// SensorClient interface for gRPC client
type SensorClient interface {
	SendSensorData(ctx context.Context, data *entities.SensorData) error
	SendSensorDataBatch(ctx context.Context, data []*entities.SensorData) error
	HealthCheck(ctx context.Context) error
	Close() error
}
