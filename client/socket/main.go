package main

import (
	"context"
	"github.com/douglasmakey/pocketknife/tracker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"net"
	"os"
	"runtime/pprof"
	"time"
)

const (
	PROTOCOL = "unix"
	SOCKET   = "/tmp/grpc.sock"

	REQUEST = 10_000
)

func main() {
	f, err := os.Create("unix.prof")
	if err != nil {
		log.Fatal(err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}

	defer pprof.StopCPUProfile()
	defer tracker.LogTimeTrack(time.Now(), "10k request to unix socket")

	dialer := func(addr string, t time.Duration) (net.Conn, error) {
		return net.Dial(PROTOCOL, addr)
	}

	conn, err := grpc.Dial(SOCKET, grpc.WithInsecure(), grpc.WithDialer(dialer))
	if err != nil {
		log.Fatal(err)
	}

	health := grpc_health_v1.NewHealthClient(conn)
	ctx := context.Background()

	for i := 0; i < REQUEST; i++ {
		heal, err := health.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%d : server status is %v", i, heal.GetStatus())
	}
}
