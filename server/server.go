package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"telnet/system"
)

const (
	DefaultPort           = 9000
	DefaultMaxConnections = 10
	DefaultBufSize        = 1024
)

var (
	ErrZeroPort    = errors.New("port value is zero")
	ErrZeroMaxConn = errors.New("max connections is zero")
)

type Server struct {
	ch      chan int8
	stopCh  chan struct{}
	wg      sync.WaitGroup
	connId  uint64
	Port    uint16
	BufSize int
}

type stopServerHandler struct {
	stopCh chan struct{}
}

func (callback *stopServerHandler) NotifyInterrupt() {
	fmt.Println("Notify server to stop properly")
	close(callback.stopCh)
}

func (s *Server) Start() error {
	fmt.Printf("Starting server at port %d\n", s.Port)
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))

	if err != nil {
		return err
	}

	for {
		select {
		case <-s.ch:
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Error accepting connection: %s\n", err)
				os.Exit(2)
			}
			s.wg.Add(1)
			go handleConnection(s.connId, conn, s.stopCh, s.ch, &s.wg)
			s.connId++

		case <-s.stopCh:
			fmt.Println("Received stop program signal...")
			s.wg.Wait()
			err := listener.Close()
			if err != nil {
				fmt.Printf("Error closing listener: %s\n", err)
				os.Exit(3)
			}
			fmt.Println("Port successfully closed")
			os.Exit(0)
		}
	}

	return nil
}

func handleConnection(id uint64, conn net.Conn, stopCh chan struct{}, doneCh chan int8, wg *sync.WaitGroup) {
	fmt.Printf("(id: %d) New connection created\n", id)
	reader := bufio.NewReader(conn)

	defer func(id uint64, wg *sync.WaitGroup) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("(id: %d) Error closing connection", id)
		}
		wg.Done()
	}(id, wg)

	var stop bool

	for {
		select {
		case <-stopCh:
			stop = true

		default:
			stop = false
		}

		if stop {
			fmt.Printf("(id: %d) Stop event received...", id)
			return
		}

		data, isPrefix, err := reader.ReadLine()

		if err != nil {
			fmt.Printf("(id: %d) Error reading from connection: %s", id, err)
			doneCh <- 1
			return
		}

		dataStr := string(data)

		if isPrefix {
			if dataStr == "exit" {
				fmt.Printf("(id: %d) Exit command received", id)
				doneCh <- 1
			} else {
				fmt.Print(string(data))
			}
		} else {
			fmt.Println(string(data))
		}
	}
}

func New(port, maxConnections uint16) (*Server, error) {
	if port == 0 {
		return nil, ErrZeroPort
	}

	if maxConnections == 0 {
		return nil, ErrZeroMaxConn
	}

	//Channel for stop notifications
	stopChannel := make(chan struct{})

	//Register callback to be processed if CTRL+C is pushed
	system.RegisterCallback(&stopServerHandler{stopCh: stopChannel})

	return &Server{
		ch:      make(chan int8, maxConnections),
		stopCh:  stopChannel,
		Port:    port,
		BufSize: DefaultBufSize,
	}, nil
}

func Start() error {
	return StartAtPort(DefaultPort)
}

func StartAtPort(port uint16) error {
	s, err := New(port, DefaultMaxConnections)
	if err != nil {
		return err
	}
	return s.Start()
}
