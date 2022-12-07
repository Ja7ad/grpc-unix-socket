package tracker

import (
	"log"
	"time"
)

func LogTimeTrack(start time.Time, name string) {
	log.Printf("%s took %s \n", name, time.Since(start))
}

func TimeTrack(start time.Time, name string) time.Duration {
	return time.Since(start)
}
