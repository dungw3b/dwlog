package dwlog

import (
	"fmt"
	"time"
	"sync"
	"strings"
	"github.com/nats-io/nats.go"
)

type Server struct {
	Subject string
	Servers []string
	Timeout time.Duration
	conn *nats.Conn
	encode *nats.EncodedConn
}

func (s *Server) Run() error {
	if err := s.Connect(); err != nil {
		return err
	}
	defer s.Close()

	wg := sync.WaitGroup{}
	wg.Add(1)

	_, err := s.encode.Subscribe(s.Subject, func(m *Message){
		fmt.Println(m)
		wg.Done()
	})
	if err != nil {
		fmt.Println(err)
	}

	wg.Wait()
	return nil
}

func (s *Server) Connect() error {
	if len(s.Subject) == 0 {
		return fmt.Errorf("subject is required")
	}
	if len(s.Servers) == 0 {
		return fmt.Errorf("list servers is required")
	}
	if s.Timeout == 0 {
		s.Timeout = 20 * time.Second
	}

	var err error
	s.conn, err = nats.Connect(strings.Join(s.Servers, ","), nats.Name("DWLog server"), nats.Timeout(s.Timeout), nats.NoEcho())
	if err != nil {
		return err
	}
	s.encode, err = nats.NewEncodedConn(s.conn, nats.JSON_ENCODER)
	if err != nil {
		s.Close()
		return fmt.Errorf("encode json error")
	}
	return nil
}

func (s *Server) Close() {
	if s.conn != nil {
		if s.encode != nil {
			s.encode.Close()
		}
		s.conn.Close()
	}
}