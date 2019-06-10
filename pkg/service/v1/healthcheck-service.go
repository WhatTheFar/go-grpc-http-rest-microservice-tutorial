package v1

import (
	"context"

	v1 "github.com/amsokol/go-grpc-http-rest-microservice-tutorial/pkg/api/v1"
)

type HealthcheckServiceServer struct {
}

func NewHealthcheckServiceServer() v1.HealthServer {
	return &HealthcheckServiceServer{}
}

func (s *HealthcheckServiceServer) Check(ctx context.Context, healthcheckReq *v1.HealthCheckRequest) (*v1.HealthCheckResponse, error) {
	healthcheck := &v1.HealthCheckResponse{
		Status: v1.HealthCheckResponse_SERVING,
	}
	return healthcheck, nil
}
