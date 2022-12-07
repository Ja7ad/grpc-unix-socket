package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	// unix socket
	PROTOCOL = "unix"
	SOCKET   = "/tmp/grpc.sock"

	// tcp protocol
	PROTOCOL_TCP = "tcp"
	ADDR         = "localhost:3000"
)

func main() {
	ln, err := net.Listen(PROTOCOL, SOCKET)
	if err != nil {
		log.Fatal(err)
	}

	tcpLn, err := net.Listen(PROTOCOL_TCP, ADDR)
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(SOCKET)
		os.Exit(1)
	}()

	srv := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	reflection.Register(srv)

	go func() {
		log.Printf("grpc ran on tcp protocol %s", ADDR)
		log.Fatal(srv.Serve(tcpLn))
	}()

	log.Printf("grpc ran on unix socket protocol %s", SOCKET)
	log.Fatal(srv.Serve(ln))
}
