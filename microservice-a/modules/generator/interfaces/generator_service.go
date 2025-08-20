package interfaces

import (
	"context"
	"time"

	"github.com/worlder-team/microservice-server/microservice-a/modules/generator/entities"
)

// GeneratorService interface
type GeneratorService interface {
	StartGeneration(ctx context.Context) error
	StopGeneration()
	SetFrequency(frequency string) error
	GetFrequency() time.Duration
	GetStatus() *entities.GeneratorStatus
	IsRunning() bool
}
