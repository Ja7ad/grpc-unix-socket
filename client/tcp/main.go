package main

import (
	"context"
	"github.com/douglasmakey/pocketknife/tracker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

const (
	PROTOCOL = "tcp"
	ADDR     = "localhost:3000"

	REQUEST = 1
)

func main() {
	f, err := os.Create("tcp2.prof")
	if err != nil {
		log.Fatal(err)
	}

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}

	defer pprof.StopCPUProfile()
	defer tracker.LogTimeTrack(time.Now(), "10k request to tcp protocol")

	conn, err := grpc.Dial(ADDR, grpc.WithInsecure())
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
