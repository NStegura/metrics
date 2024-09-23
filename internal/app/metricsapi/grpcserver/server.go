package grpcserver

import (
	"context"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/NStegura/metrics/config"
	blModels "github.com/NStegura/metrics/internal/business/models"
	pb "github.com/NStegura/metrics/pkg/api"
)

type MetricsGRPCServer struct {
	pb.UnimplementedMetricsApiServer

	cfg *config.SrvConfig
	bll Bll

	logger *logrus.Logger
}

func New(cfg *config.SrvConfig, bll Bll, logger *logrus.Logger) (*MetricsGRPCServer, error) {
	return &MetricsGRPCServer{
		cfg:    cfg,
		bll:    bll,
		logger: logger,
	}, nil
}

func (s *MetricsGRPCServer) Start(opts ...grpc.ServerOption) error {
	s.logger.Infof("starting GRPCServer %s", s.cfg.GrpcAddr)
	lis, err := net.Listen("tcp", s.cfg.GrpcAddr)
	defer func() {
		if err = lis.Close(); err != nil {
			s.logger.Error("failed to close listener")
		}
	}()
	if err != nil {
		return fmt.Errorf("failed to create network listener: %w", err)
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMetricsApiServer(grpcServer, s)
	if err = grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}
	return nil
}

func (s *MetricsGRPCServer) UpdateAllMetrics(ctx context.Context, req *pb.MetricsList) (*pb.UpdateResponse, error) {
	for _, metric := range req.Metrics {
		switch metric.Mtype {
		case pb.MetricType_GAUGE:
			err := s.bll.UpdateGaugeMetric(ctx, blModels.GaugeMetric{
				Name:  metric.Id,
				Type:  "gauge",
				Value: metric.Value,
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update gauge metric")
			}
		case pb.MetricType_COUNTER:
			err := s.bll.UpdateCounterMetric(ctx, blModels.CounterMetric{
				Name:  metric.Id,
				Type:  "counter",
				Value: metric.Delta,
			})
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to update counter metric")
			}
		default:
			return nil, status.Errorf(codes.InvalidArgument, "unknown metric type")
		}
	}
	return &pb.UpdateResponse{Message: "Metrics updated successfully"}, nil
}

func (s *MetricsGRPCServer) GetPing(_ context.Context, _ *emptypb.Empty) (*pb.Pong, error) {
	return &pb.Pong{
		Pong: true,
	}, nil
}
