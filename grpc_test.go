package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

func Benchmark_UNIX(b *testing.B) {
	f, err := os.Create("unix.prof")
	if err != nil {
		b.Fatal(err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		b.Fatal(err)
	}

	defer pprof.StopCPUProfile()

	dialer := func(addr string, t time.Duration) (net.Conn, error) {
		return net.Dial(PROTOCOL, addr)
	}

	conn, err := grpc.Dial(SOCKET, grpc.WithInsecure(), grpc.WithDialer(dialer))
	if err != nil {
		b.Fatal(err)
	}

	health := grpc_health_v1.NewHealthClient(conn)
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		_, err := health.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_TCP(b *testing.B) {
	f, err := os.Create("tcp.prof")
	if err != nil {
		b.Fatal(err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		b.Fatal(err)
	}

	defer pprof.StopCPUProfile()

	conn, err := grpc.Dial(ADDR, grpc.WithInsecure())
	if err != nil {
		b.Fatal(err)
	}

	health := grpc_health_v1.NewHealthClient(conn)
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		_, err := health.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			b.Fatal(err)
		}
	}
}
