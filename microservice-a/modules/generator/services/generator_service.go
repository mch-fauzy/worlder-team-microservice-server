package services

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/worlder-team/microservice-server/microservice-a/modules/generator/entities"
	"github.com/worlder-team/microservice-server/microservice-a/modules/generator/interfaces"
	"github.com/worlder-team/microservice-server/shared/constants"
	"github.com/worlder-team/microservice-server/shared/utils"
)

type generatorService struct {
	grpcClient    interfaces.SensorClient
	sensorType    string
	frequency     time.Duration
	isRunning     bool
	mu            sync.RWMutex
	stopChan      chan struct{}
	ctx           context.Context
	cancel        context.CancelFunc
	lastGenerated time.Time
	totalSent     int64
	errors        int64
}

// NewGeneratorService creates a new generator service
func NewGeneratorService(grpcClient interfaces.SensorClient, sensorType, frequency string) interfaces.GeneratorService {
	freq, err := utils.ParseDuration(frequency)
	if err != nil {
		freq = time.Second // Default to 1 second
	}

	return &generatorService{
		grpcClient: grpcClient,
		sensorType: sensorType,
		frequency:  freq,
		stopChan:   make(chan struct{}),
	}
}

// StartGeneration starts generating sensor data
func (s *generatorService) StartGeneration(ctx context.Context) error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return fmt.Errorf("generator is already running")
	}
	s.isRunning = true
	s.stopChan = make(chan struct{})

	// Create a new context that's independent of the request context
	s.ctx, s.cancel = context.WithCancel(context.Background())
	frequency := s.frequency
	s.mu.Unlock()

	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			s.mu.Lock()
			s.isRunning = false
			s.mu.Unlock()
			return s.ctx.Err()
		case <-s.stopChan:
			s.mu.Lock()
			s.isRunning = false
			s.mu.Unlock()
			return nil
		case <-ticker.C:
			s.generateAndSend(s.ctx)
		}
	}
}

// StopGeneration stops the data generation
func (s *generatorService) StopGeneration() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		s.isRunning = false
		if s.cancel != nil {
			s.cancel()
		}
		close(s.stopChan)
	}
}

// SetFrequency sets the generation frequency
func (s *generatorService) SetFrequency(frequency string) error {
	freq, err := utils.ParseDuration(frequency)
	if err != nil {
		return fmt.Errorf("invalid frequency format: %v", err)
	}

	s.mu.Lock()
	wasRunning := s.isRunning
	s.frequency = freq
	s.mu.Unlock()

	// If generator is running, restart it automatically with new frequency
	if wasRunning {
		s.StopGeneration()
		// Give it a moment to fully stop
		time.Sleep(100 * time.Millisecond)
		// Start with new frequency in background
		go s.StartGeneration(context.Background())
	}

	return nil
}

// GetFrequency returns the current frequency
func (s *generatorService) GetFrequency() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.frequency
}

// GetStatus returns the current generator status
func (s *generatorService) GetStatus() *entities.GeneratorStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &entities.GeneratorStatus{
		IsRunning:     s.isRunning,
		SensorType:    s.sensorType,
		Frequency:     s.frequency,
		LastGenerated: s.lastGenerated,
		TotalSent:     s.totalSent,
		Errors:        s.errors,
	}
}

// IsRunning returns whether the generator is currently running
func (s *generatorService) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}

// generateAndSend generates and sends sensor data
func (s *generatorService) generateAndSend(ctx context.Context) {
	data := s.generateSensorData()

	if err := s.grpcClient.SendSensorData(ctx, data); err != nil {
		s.mu.Lock()
		s.errors++
		s.mu.Unlock()
		// Log error but continue generation
		fmt.Printf("Error sending sensor data: %v\n", err)
		return
	}

	s.mu.Lock()
	s.totalSent++
	s.lastGenerated = time.Now()
	s.mu.Unlock()
}

// generateSensorData generates sensor data based on sensor type
func (s *generatorService) generateSensorData() *entities.SensorData {
	now := time.Now()

	// Generate ID1 (random uppercase string)
	id1 := utils.GenerateID(4) // 8 character hex string, uppercase

	// Generate ID2 (random integer)
	id2 := rand.Int31n(10000) // Random number between 0-9999

	// Generate sensor value based on type
	var value float64
	switch s.sensorType {
	case constants.SensorTypeTemperature:
		// Temperature: -10 to 50 degrees Celsius
		value = -10 + rand.Float64()*60
	case constants.SensorTypeHumidity:
		// Humidity: 0 to 100 percent
		value = rand.Float64() * 100
	case constants.SensorTypePressure:
		// Pressure: 980 to 1030 hPa
		value = 980 + rand.Float64()*50
	case constants.SensorTypeLight:
		// Light: 0 to 1000 lux
		value = rand.Float64() * 1000
	case constants.SensorTypeMotion:
		// Motion: 0 or 1 (binary)
		value = float64(rand.Intn(2))
	default:
		// Default: random value between 0-100
		value = rand.Float64() * 100
	}

	return &entities.SensorData{
		SensorValue: value,
		SensorType:  s.sensorType,
		ID1:         id1,
		ID2:         id2,
		Timestamp:   now,
	}
}
