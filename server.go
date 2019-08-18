package dwlog

import (
	"net"
)

type Server struct {

}

func (s *Server) GetMessage(ctx context.Context, m *Message) (*Response, error) {
	return &Response{result: true}
}

func (s *Server) Run() error {
	listen, err := net.Listen("tcp", 2701)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	RegisterLogServiceServer(grpcServer, &Server{})
	grpcServer.Serve(listen)
	return nil
}