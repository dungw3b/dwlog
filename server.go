package dwlog

import (
	"os"
	"fmt"
	"net"
	"time"
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

	mux *sync.Mutex
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
		return &MessageResponse{Result: false}, nil
	}

	return &MessageResponse{Result: true}, nil
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

func (s *Server) Run() error {
	var (
		hostname, port, info string
	)
	if len(s.Listen) == 0 {
		return fmt.Errorf("service.listen is required")
	}
	if s.Port == 0 {
		return fmt.Errorf("service.port is required")
	}
	if len(s.Data) == 0 {
		s.Data = "./"
	}
	s.mux = &sync.Mutex{}

	port = strconv.Itoa(int(s.Port))
	listen, err := net.Listen("tcp", s.Listen +":"+ port)
	if err != nil {
		return err
	}

	hostname,_ = os.Hostname()
	info = "["+ time.Now().Format("2006-01-02 15:04:05") +"]"
	info = info + "["+ hostname +"]"
	info = info + "[\033[32mINFO\033[0m]"
	info = info + " Service dwlog started on "+ s.Listen +":"+ port
	fmt.Println(info)

	grpcServer := grpc.NewServer()
	RegisterLogServiceServer(grpcServer, s)
	grpcServer.Serve(listen)
	return nil
}