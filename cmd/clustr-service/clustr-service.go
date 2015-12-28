package main

import (
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/dmage/clustr/watcher"
)

var (
	serviceUnitFile = kingpin.Flag("unit", "Service unit configuration file.").Short('u').Required().String()
	start           = kingpin.Command("start", "Start service.")
	stop            = kingpin.Command("stop", "Stop service.")
	status          = kingpin.Command("status", "Check if service is running.")
)

func main() {
	command := kingpin.Parse()

	p := watcher.PrisonerFromFile(*serviceUnitFile)
	switch command {
	case start.FullCommand():
		p.Start()
	case stop.FullCommand():
		p.Stop()
	case status.FullCommand():
		if p.IsRunning() {
			log.Print(p.Name, " is running")
		} else {
			log.Print(p.Name, " stopped")
			os.Exit(1)
		}
	default:
		panic("unexpected command: " + command)
	}
}
