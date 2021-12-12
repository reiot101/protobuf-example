package grpcinfo

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	dpb "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/prasek/protoer/proto"
	"google.golang.org/grpc"
)

type Registry interface {
	Load(srv *grpc.Server) error
	LoadFile(file string) error
	Server(fqn string) Server
}

func NewRegistry() Registry {
	return &registry{
		servers: make(map[string]*server),
	}
}

type registry struct {
	servers map[string]*server
}

func (r *registry) Load(srv *grpc.Server) error {
	for name, info := range srv.GetServiceInfo() {
		file, ok := info.Metadata.(string)
		if !ok {
			return fmt.Errorf("Service %q has unexpected metadata. Expecting a string, got %v", name, info.Metadata)
		}
		if file == "reflection/grpc_reflection_v1alpha/reflection.proto" {
			continue
		}
		fmt.Println("file:", file)
		err := r.LoadFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *registry) LoadFile(file string) error {
	fd, err := loadFileDescriptorProto(file)
	if err != nil {
		return err
	}

	pkg := fd.GetPackage()

	merge := func(a, b string) string {
		if a == "" {
			return b
		} else {
			return a + "." + b
		}
	}

	for i := range fd.Service {
		svc := fd.Service[i]
		fqn := merge(pkg, svc.GetName())
		for j := range svc.Method {
			m := svc.Method[j]
			fqnMethod := fmt.Sprintf("/%s/%s", fqn, m.GetName())
			r.servers[fqnMethod] = &server{service: &serviceExtention{proto: svc}, method: &methodExtention{proto: m}}
		}
	}
	return nil
}

func (r *registry) Server(fqn string) Server {
	return r.servers[fqn]
}

// loadFileDescriptor loads a registered descriptor and decodes it. If the given
// name cannot be loaded but is a known standard name, an alias will be tried by the proto,
// so the standard files can be loaded even if linked against older "known bad"
// versions of packages.
func loadFileDescriptorProto(file string) (*dpb.FileDescriptorProto, error) {
	fdb := proto.FileDescriptor(file)
	if fdb == nil {
		return nil, fmt.Errorf("Missing file descriptor %s.", file)
	}

	fd, err := decodeFileDescriptorProto(file, fdb)
	if err != nil {
		return nil, err
	}

	// the file descriptor may have been laoded with an alias,
	// so we ensure the specified name to ensure it can be linked.
	fd.Name = proto.String(file)

	return fd, nil
}

// decodeFileDescriptorProto decodes the bytes of a registered file descriptor.
// Registered file descriptors are first "proto encoded" (e.g. binary format
// for the descriptor protos) and then gzipped. So this function gunzips and
// then unmarshals into a descriptor proto.
func decodeFileDescriptorProto(element string, fdb []byte) (*dpb.FileDescriptorProto, error) {
	raw, err := decompress(fdb)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress %q descriptor: %v", element, err)
	}
	fd := dpb.FileDescriptorProto{}
	if err := proto.Unmarshal(raw, &fd); err != nil {
		return nil, fmt.Errorf("bad descriptor for %q: %v", element, err)
	}
	return &fd, nil
}

func decompress(b []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("bad gzipped descriptor: %v", err)
	}
	out, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("bad gzipped descriptor: %v", err)
	}
	return out, nil
}
