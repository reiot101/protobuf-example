package server

import (
	"context"
	"log"

	"github.com/reiot777/protobuf-example/internal/grpcinfo"
	"github.com/reiot777/protobuf-example/packet"
	"google.golang.org/grpc"
)

func (s *Server) EveryAuth(reg grpcinfo.Registry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		log.Println("[auth] request ", req)

		mi := reg.Server(info.FullMethod)
		if v, err := mi.Service().Get(packet.E_OauthScopes); err != nil {
			log.Println("service err: ", err)
		} else {
			log.Println("service: ", v)
		}
		if v, err := mi.Method().Get(packet.E_MethodSignature); err != nil {
			log.Println("method err: ", err)
		} else {
			log.Println("method: ", v)
		}

		resp, err = handler(ctx, req)
		log.Println("[auth] response ", req)
		return
	}
}
