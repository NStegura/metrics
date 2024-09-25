package metric

import (
	"context"
	"fmt"
	"time"

	"github.com/NStegura/metrics/internal/utils/ip"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/NStegura/metrics/internal/clients/base"
	"github.com/NStegura/metrics/pkg/api"
)

type GRPCClient struct {
	*base.BaseClient
	conn   *grpc.ClientConn
	client api.MetricsApiClient // это место для gRPC клиента
}

func NewGRPCClient(addr string, options ...base.Option) (*GRPCClient, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	bc, err := base.NewBaseClient(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to init client: %w", err)
	}

	return &GRPCClient{
		BaseClient: bc,
		conn:       conn,
		client:     api.NewMetricsApiClient(conn), // инициализация клиента
	}, nil
}

// UpdateMetrics обновляет набор метрик.
func (c *GRPCClient) UpdateMetrics(ctx context.Context, metrics []Metrics) error {
	ml := c.convert(metrics)

	ctx, err := c.prepareCtx(ctx)
	if err != nil {
		return fmt.Errorf("failed to prepare ctx: %w", err)
	}
	_, err = c.Execute(
		func() (any, error) {
			return c.client.UpdateAllMetrics(ctx, &api.MetricsList{Metrics: ml}) //nolint:wrapcheck // proxy
		},
		c.conn.Target(),
		"grpc",
	)
	if err != nil {
		return fmt.Errorf("failed to send metrics via gRPC: %w", err)
	}

	return nil
}

func (c *GRPCClient) prepareCtx(ctx context.Context) (context.Context, error) {
	selfIP, err := ip.GetIP()
	if err != nil {
		return ctx, fmt.Errorf("failed to get ip: %w", err)
	}

	ctx = metadata.AppendToOutgoingContext(ctx,
		"when", time.Now().Format(time.RFC3339),
		"sender", "agent",
		"ip", selfIP,
	)
	return ctx, nil
}

func (c *GRPCClient) convert(metrics []Metrics) []*api.Metric {
	var mtype api.MetricType

	ml := make([]*api.Metric, 0, len(metrics))
	for _, m := range metrics {
		if m.MType == "gauge" {
			mtype = api.MetricType_GAUGE
		} else if m.MType == "counter" {
			mtype = api.MetricType_COUNTER
		}
		metric := api.Metric{
			Id:    m.ID,
			Mtype: mtype,
		}
		if m.Value != nil {
			metric.Value = *m.Value
		}
		if m.Delta != nil {
			metric.Delta = *m.Delta
		}
		ml = append(ml, &metric)
	}
	return ml
}
