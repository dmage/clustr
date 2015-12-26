package main

import (
	"log"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/dmage/clustr/service"
	"github.com/dmage/clustr/unit"
)

var (
	serviceUnitFile = kingpin.Flag("unit", "Service unit configuration file.").Short('u').Required().String()
	start           = kingpin.Command("start", "Start service.")
	stop            = kingpin.Command("stop", "Stop service.")
)

func main() {
	command := kingpin.Parse()

	serviceUnit, err := unit.ServiceUnitFromFile(*serviceUnitFile)
	if err != nil {
		log.Fatal("failed to load service: ", err)
	}

	service := service.NewService("sleep.service", &serviceUnit.Service)

	switch command {
	case start.FullCommand():
		err = service.Start()
		if err != nil {
			log.Fatal("failed to start service: ", err)
		}
	case stop.FullCommand():
		err = service.Stop()
		if err != nil {
			log.Fatal("failed to stop service: ", err)
		}
	default:
		panic("unexpected command: " + command)
	}
}
