package grpcinfo

type Server interface {
	Service() Extention
	Method() Extention
}

type server struct {
	service *serviceExtention
	method  *methodExtention
}

func (s *server) Service() Extention {
	return s.service
}

func (s *server) Method() Extention {
	return s.method
}
