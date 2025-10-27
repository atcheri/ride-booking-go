package grpcclient

import (
	"github.com/atcheri/ride-booking-go/shared/tracing"
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/driver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type driverServiceClient struct {
	Client pb.DriverServiceClient
	conn   *grpc.ClientConn
}

func NewDriverServiceClient(url string) (*driverServiceClient, error) {
	dialOptions := append(
		tracing.DialOptionsWithTracing(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	conn, err := grpc.NewClient(url, dialOptions...)
	if err != nil {
		return nil, err
	}

	client := pb.NewDriverServiceClient(conn)
	return &driverServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *driverServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
