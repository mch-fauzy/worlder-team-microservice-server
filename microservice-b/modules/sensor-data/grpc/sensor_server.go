package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/entities"
	"github.com/worlder-team/microservice-server/microservice-b/modules/sensor-data/interfaces"
	pb "github.com/worlder-team/microservice-server/shared/proto/sensor"
)

type sensorServer struct {
	pb.UnimplementedSensorServiceServer
	sensorService interfaces.SensorServiceInterface
}

// NewSensorServer creates a new gRPC sensor server
func NewSensorServer(sensorService interfaces.SensorServiceInterface) *sensorServer {
	return &sensorServer{
		sensorService: sensorService,
	}
}

// RegisterSensorServer registers the sensor server with gRPC
func RegisterSensorServer(s *grpc.Server, sensorService interfaces.SensorServiceInterface) {
	pb.RegisterSensorServiceServer(s, NewSensorServer(sensorService))
}

// SendSensorData handles single sensor data reception
func (s *sensorServer) SendSensorData(ctx context.Context, req *pb.SensorData) (*pb.SensorResponse, error) {
	// Convert protobuf to domain entity
	sensorData := &entities.SensorData{
		SensorValue: req.SensorValue,
		SensorType:  req.SensorType,
		ID1:         req.Id1,
		ID2:         req.Id2,
		Timestamp:   req.Timestamp.AsTime(),
	}

	// Save to database
	err := s.sensorService.CreateSensorData(ctx, sensorData)
	if err != nil {
		return &pb.SensorResponse{
			Success: false,
			Message: "Failed to save sensor data",
			Error:   err.Error(),
		}, nil
	}

	return &pb.SensorResponse{
		Success: true,
		Message: "Sensor data saved successfully",
	}, nil
}

// SendSensorDataBatch handles batch sensor data reception
func (s *sensorServer) SendSensorDataBatch(ctx context.Context, req *pb.SensorDataBatch) (*pb.SensorResponse, error) {
	// Convert protobuf batch to domain entities
	var sensorDataBatch []*entities.SensorData
	for _, data := range req.Data {
		sensorData := &entities.SensorData{
			SensorValue: data.SensorValue,
			SensorType:  data.SensorType,
			ID1:         data.Id1,
			ID2:         data.Id2,
			Timestamp:   data.Timestamp.AsTime(),
		}
		sensorDataBatch = append(sensorDataBatch, sensorData)
	}

	// Save batch to database
	err := s.sensorService.CreateSensorDataBatch(ctx, sensorDataBatch)
	if err != nil {
		return &pb.SensorResponse{
			Success: false,
			Message: "Failed to save sensor data batch",
			Error:   err.Error(),
		}, nil
	}

	return &pb.SensorResponse{
		Success: true,
		Message: "Sensor data batch saved successfully",
	}, nil
}

// HealthCheck handles health check requests
func (s *sensorServer) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}

// Helper function to convert time to protobuf timestamp
func timeToTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
