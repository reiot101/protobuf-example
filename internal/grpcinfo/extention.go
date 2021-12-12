package grpcinfo

import (
	dpb "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/prasek/protoer/proto"
)

var (
	_ Extention = &serviceExtention{}
	_ Extention = &methodExtention{}
)

type Extention interface {
	Get(ext interface{}) (interface{}, error)
	Has(ext interface{}, ifnotset bool) bool
}

type serviceExtention struct {
	proto *dpb.ServiceDescriptorProto
}

func (s *serviceExtention) Get(ext interface{}) (interface{}, error) {
	return proto.GetExtension(s.proto.GetOptions(), ext)
}

func (s *serviceExtention) Has(ext interface{}, ifnotset bool) bool {
	return proto.GetBoolExtension(s.proto.GetOptions(), ext, ifnotset)
}

type methodExtention struct {
	proto *dpb.MethodDescriptorProto
}

func (s *methodExtention) Get(ext interface{}) (interface{}, error) {
	return proto.GetExtension(s.proto.GetOptions(), ext)
}

func (s *methodExtention) Has(ext interface{}, ifnotset bool) bool {
	return proto.GetBoolExtension(s.proto.GetOptions(), ext, ifnotset)
}
