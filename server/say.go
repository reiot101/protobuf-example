package server

import (
	"context"
	"time"

	"github.com/reiot777/protobuf-example/packet"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ packet.SayServer = &Service{}

func (s *Service) Hello(ctx context.Context, req *packet.HelloRequest) (*packet.HelloResponse, error) {
	return &packet.HelloResponse{
		Msg: "Hello " + req.Msg,
	}, nil
}

func (s *Service) Ping(ctx context.Context, _ *emptypb.Empty) (*packet.PingResponse, error) {
	return &packet.PingResponse{
		Ts: time.Now().UnixNano(),
	}, nil
}
