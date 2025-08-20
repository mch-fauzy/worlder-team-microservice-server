package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/worlder-team/microservice-server/microservice-a/modules/generator/entities"
	"github.com/worlder-team/microservice-server/microservice-a/modules/generator/interfaces"
	pb "github.com/worlder-team/microservice-server/shared/proto/sensor"
)

type sensorClient struct {
	conn   *grpc.ClientConn
	client pb.SensorServiceClient
}

// NewSensorClient creates a new gRPC sensor client
func NewSensorClient(serverAddress string) (interfaces.SensorClient, error) {
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %v", err)
	}

	client := pb.NewSensorServiceClient(conn)

	return &sensorClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the gRPC connection
func (c *sensorClient) Close() error {
	return c.conn.Close()
}

// SendSensorData sends a single sensor data to the server
func (c *sensorClient) SendSensorData(ctx context.Context, data *entities.SensorData) error {
	pbData := &pb.SensorData{
		SensorValue: data.SensorValue,
		SensorType:  data.SensorType,
		Id1:         data.ID1,
		Id2:         data.ID2,
		Timestamp:   timestamppb.New(data.Timestamp),
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	response, err := c.client.SendSensorData(ctx, pbData)
	if err != nil {
		return fmt.Errorf("failed to send sensor data: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("server error: %s", response.Error)
	}

	return nil
}

// SendSensorDataBatch sends a batch of sensor data to the server
func (c *sensorClient) SendSensorDataBatch(ctx context.Context, data []*entities.SensorData) error {
	var pbDataBatch []*pb.SensorData

	for _, sensorData := range data {
		pbData := &pb.SensorData{
			SensorValue: sensorData.SensorValue,
			SensorType:  sensorData.SensorType,
			Id1:         sensorData.ID1,
			Id2:         sensorData.ID2,
			Timestamp:   timestamppb.New(sensorData.Timestamp),
		}
		pbDataBatch = append(pbDataBatch, pbData)
	}

	request := &pb.SensorDataBatch{
		Data: pbDataBatch,
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	response, err := c.client.SendSensorDataBatch(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to send sensor data batch: %v", err)
	}

	if !response.Success {
		return fmt.Errorf("server error: %s", response.Error)
	}

	return nil
}

// HealthCheck checks the health of the gRPC server
func (c *sensorClient) HealthCheck(ctx context.Context) error {
	request := &pb.HealthCheckRequest{
		Service: "microservice-a",
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	response, err := c.client.HealthCheck(ctx, request)
	if err != nil {
		return fmt.Errorf("health check failed: %v", err)
	}

	if response.Status != pb.HealthCheckResponse_SERVING {
		return fmt.Errorf("server not serving")
	}

	return nil
}
