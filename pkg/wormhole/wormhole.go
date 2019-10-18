package wormhole

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Wormhole struct {
	listener   net.Listener
	connection *grpc.ClientConn
}

func (w *Wormhole) RunServer(port string, registerServer func(*grpc.Server)) error {
	var err error
	w.listener, err = net.Listen("tcp", port)
	if err != nil {
		return err
	}
	GrpcServer := grpc.NewServer()

	registerServer(GrpcServer)
	quitChan := make(chan error)
	go func(out chan<- error) {
		fmt.Println("Start server")
		err := GrpcServer.Serve(w.listener)
		out <- err
	}(quitChan)

	res := <-quitChan
	fmt.Println("Quit server")
	return res
}

func (w *Wormhole) RunClient(host, port string, registerClient func(*grpc.ClientConn)) error {
	var err error
	w.connection, err = grpc.Dial(host+"+"+port, grpc.WithInsecure())
	if err != nil {
		// fmt.Println("Cannot connect to service")
		return err
	}
	registerClient(w.connection)
	return nil
}
