package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"telnet/system"
	"time"
)

const defaultTimeout = time.Second * 15

type closeClientHandler struct {
	connection net.Conn
}

func (handler *closeClientHandler) NotifyInterrupt() {
	fmt.Println("Closing connection")
	err := handler.connection.Close()
	if err != nil {
		fmt.Printf("Error closing connection: %s\n", err)
	} else {
		fmt.Println("Connection closed")
	}
	os.Exit(0)
}

func Start(host string, port uint16) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), defaultTimeout)

	if err != nil {
		return err
	}

	fmt.Printf("Connected to %s:%d\n", host, port)

	//Register CTRL+C handler
	interruptHandler := &closeClientHandler{connection: conn}
	system.RegisterCallback(interruptHandler)

	reader := bufio.NewReader(os.Stdin)

	for {
		data, _, _ := reader.ReadLine()

		if string(data) == "exit" {
			//Ignore errors
			_, _ = conn.Write([]byte("exit"))
			interruptHandler.NotifyInterrupt()
			return nil
		}

		data = append(data, '\n')
		n, err := conn.Write(data)

		if err != nil {
			interruptHandler.NotifyInterrupt()
			return errors.New(fmt.Sprintf("Error sending data: %s", err))
		}

		if n < len(data) {
			sended := n

			for {
				n, err = conn.Write(data[sended:])

				if err != nil {
					interruptHandler.NotifyInterrupt()
					return errors.New(fmt.Sprintf("Error sending data: %s", err))
				}

				sended += n

				if sended == len(data) {
					break
				}
			}
		}
	}
}
