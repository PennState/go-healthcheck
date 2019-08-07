package cpu

import (
	healthcheck "github.com/PennState/go-healthcheck/pkg/health"
	linuxproc "github.com/c9s/goprocinfo/linux"
	log "github.com/sirupsen/logrus"
)

type CPUCheck struct {
}

//See: https://github.com/c9s/goprocinfo
//See: https://www.linuxhowtos.org/System/procstat.htm
func Check() ([]healthcheck.Check, healthcheck.Status) {
	var checks []healthcheck.Check

	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		return checks, healthcheck.Warn
	}

	log.Info(stat.CPUStats[0])

	for _, s := range stat.CPUStats {
		log.Info("-----")
		log.Info("Id: ", s.Id)
		log.Info("User: ", s.User)
		log.Info("Nice: ", s.Nice)
		log.Info("System: ", s.System)
		log.Info("Idle: ", s.Idle)
		log.Info("IOWait: ", s.IOWait)
		log.Info("IRQ: ", s.IRQ)
		log.Info("Soft IRQ: ", s.SoftIRQ)
		log.Info("Steal: ", s.Steal)
		log.Info("Guest: ", s.Guest)
		log.Info("Guest nice: ", s.GuestNice)
	}
	log.Info("-----")

	return checks, healthcheck.Pass
}
