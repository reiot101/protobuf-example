package server

import "github.com/reiot777/protobuf-example/packet"

type Service struct {
	packet.UnimplementedSayServer
}
