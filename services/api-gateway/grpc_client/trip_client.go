package grpcclient

import (
	"github.com/atcheri/ride-booking-go/shared/tracing"
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type tripServiceClient struct {
	Client pb.TripServiceClient
	conn   *grpc.ClientConn
}

func NewTripServiceClient(url string) (*tripServiceClient, error) {
	dialOptions := append(
		tracing.DialOptionsWithTracing(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	conn, err := grpc.NewClient(url, dialOptions...)
	if err != nil {
		return nil, err
	}

	client := pb.NewTripServiceClient(conn)
	return &tripServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *tripServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
