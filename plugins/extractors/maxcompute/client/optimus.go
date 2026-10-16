package client

import (
	"context"
	"fmt"
	"time"

	og "github.com/goto/meteor/metrics/otelgrpc"
	pb "github.com/goto/meteor/models/gotocompany/optimus/core/v1beta1"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	OptimusGRPCMaxClientSendSizeMB = 45
	OptimusGRPCMaxClientRecvSizeMB = 45
	optimusGRPCMaxRetry            = 3
)

//go:generate mockery --name=OptimusClient --case underscore --with-expecter --structname OptimusClientMock --filename optimus_client_mock.go --output=../mocks

type OptimusClient interface {
	ListResourceSpecification(ctx context.Context, in *pb.ListResourceSpecificationRequest, opts ...grpc.CallOption) (*pb.ListResourceSpecificationResponse, error)
	Connect(ctx context.Context, host string, maxSizeInMB int) error
	Close() error
}

type optimusClient struct {
	pb.ResourceServiceClient
	conn *grpc.ClientConn
}

func NewOptimusClient() OptimusClient {
	return &optimusClient{}
}

func (c *optimusClient) Connect(ctx context.Context, host string, maxSizeInMB int) error {
	dialTimeoutCtx, dialCancel := context.WithTimeout(ctx, 2*time.Second)
	defer dialCancel()

	conn, err := dialOptimusGRPC(dialTimeoutCtx, host, maxSizeInMB)
	if err != nil {
		return fmt.Errorf("create connection: %w", err)
	}

	c.conn = conn
	c.ResourceServiceClient = pb.NewResourceServiceClient(conn)
	return nil
}

func (c *optimusClient) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func dialOptimusGRPC(ctx context.Context, host string, maxSizeInMB int) (*grpc.ClientConn, error) {
	m := og.NewOtelGRPCMonitor(host)

	if maxSizeInMB <= 0 {
		maxSizeInMB = OptimusGRPCMaxClientRecvSizeMB
	}

	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(optimusGRPCMaxRetry),
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(OptimusGRPCMaxClientSendSizeMB<<20),
			grpc.MaxCallRecvMsgSize(maxSizeInMB<<20),
		),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_retry.UnaryClientInterceptor(retryOpts...),
			otelgrpc.UnaryClientInterceptor(),
			grpc_prometheus.UnaryClientInterceptor,
			m.UnaryClientInterceptor(),
		)),
		grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
			otelgrpc.StreamClientInterceptor(),
			grpc_prometheus.StreamClientInterceptor,
			m.StreamClientInterceptor(),
		)),
	}

	return grpc.DialContext(ctx, host, opts...)
}
