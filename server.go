/*
* @project  dwlog
* @file     server.go
* @author   dungw3b
* @date     2019-08-20
*/
package dwlog

import (
	"os"
	"fmt"
	"net"
	"sync"
	"context"
	"strconv"
	"google.golang.org/grpc"
)

type Server struct {
	Listen string
	Port uint32
	FileCount uint32
	Data string

	mux sync.Mutex
	grpcServer *grpc.Server
}

func (s *Server) Log(ctx context.Context, m *MessageRequest) (*MessageResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	line := s.CSVFormat(m)
	line = append(line, '\n')

	path := s.Data +"/"+ m.Name
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0644)
	}

	writer := &FileWriter {
		Name:  path +"/"+ m.Name +".log",
		MaxCount: int(s.FileCount),
	}
	_, err := writer.Write(line)
	if err != nil {
		fmt.Println("\033[31mWrite data into file error\033[0m", err.Error())
		return &MessageResponse{}, nil
	}

	return &MessageResponse{}, nil
}

func (s *Server) TextFormat(m *MessageRequest) []byte {
	text := "["+ m.Time +"]"+
		"["+ m.Host +"]"+
		"["+ LevelString(m.Level) +"]"+
		" "+ m.Message
	return []byte(text)
}

func (s *Server) CSVFormat(m *MessageRequest) []byte {
	var csv string
	csv = m.Time +","+ m.Host +","+ LevelString(m.Level) +","+ m.Message
	return []byte(csv)
}

func (s *Server) Run(exit chan bool) {
	var port string
	if len(s.Listen) == 0 {
		fmt.Println("service.listen is required")
		exit <- true
	}
	if s.Port == 0 {
		fmt.Println("service.port is required")
		exit <- true
	}
	if len(s.Data) == 0 {
		s.Data = "./"
	}
	port = strconv.Itoa(int(s.Port))
	listen, err := net.Listen("tcp", s.Listen +":"+ port)
	if err != nil {
		fmt.Println(err)
		exit <- true
	}

	s.grpcServer = grpc.NewServer()
	RegisterLogServiceServer(s.grpcServer, s)
	logConsole(INFO, "Service dwlog start on", s.Listen+":"+port)
	if err := s.grpcServer.Serve(listen); err != nil {
		fmt.Println(err)
		exit <- true
	}
}

func (s *Server) Close() {
	s.grpcServer.GracefulStop()
	logConsole(INFO, "Graceful stop service dwlog")
}