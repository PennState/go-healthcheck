package main

import (
	"time"

	"github.com/PennState/go-healthcheck/pkg/checks/cpu"
	log "github.com/sirupsen/logrus"
)

func main() {
	cpu := cpu.CPUCheck{}
	start := time.Now().UnixNano()
	cpu.Check()
	cpu.Check()
	end := time.Now().UnixNano()
	log.Info("Elapsed: ", end-start)
}
