package server

import (
	"context"
	"log"
	"net"
	"strconv"

	"github.com/prasek/protoer/proto"
	protoer "github.com/prasek/protoer/proto/golang"
	"github.com/reiot777/protobuf-example/internal/grpcinfo"
	"github.com/reiot777/protobuf-example/packet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Host string
	Port int
}

func (s *Server) Serve(ctx context.Context) {
	svc := Service{}

	proto.SetProtoer(protoer.NewProtoer(nil))
	info := grpcinfo.NewRegistry()
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(s.EveryAuth(info)),
	)

	go func() {
		defer srv.GracefulStop()
		<-ctx.Done()
	}()

	packet.RegisterSayServer(srv, &svc)
	reflection.Register(srv)

	if err := info.Load(srv); err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", net.JoinHostPort(s.Host, strconv.Itoa(s.Port)))
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	log.Println("gRPC server listening at ", lis.Addr())

	if err := srv.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
