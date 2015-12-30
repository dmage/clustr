package main

import (
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/dmage/clustr/highservice"
)

var (
	serviceUnitFile = kingpin.Flag("unit", "Service unit configuration file.").Short('u').Required().String()
	start           = kingpin.Command("start", "Start service.")
	stop            = kingpin.Command("stop", "Stop service.")
	status          = kingpin.Command("status", "Check if service is running.")
)

func main() {
	command := kingpin.Parse()

	s := highservice.HighServiceFromFile(*serviceUnitFile)
	switch command {
	case start.FullCommand():
		s.Start()
	case stop.FullCommand():
		s.Stop()
	case status.FullCommand():
		if s.IsRunning() {
			log.Print(s.Name, " is running")
		} else {
			log.Print(s.Name, " stopped")
			os.Exit(1)
		}
	default:
		panic("unexpected command: " + command)
	}
}
