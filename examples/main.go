package main

import (
	"time"

	"github.com/selesy/go-healthcheck/checks/cpu"
	log "github.com/sirupsen/logrus"
)

func main() {
	start := time.Now().UnixNano()
	cpu.Check()
	cpu.Check()
	end := time.Now().UnixNano()
	log.Info("Elapsed: ", end-start)
}
