package main

import (
	"log"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/dmage/clustr/logging"
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

	serviceName, serviceUnit, err := unit.ServiceUnitFromFile(*serviceUnitFile)
	if err != nil {
		log.Fatal("failed to load service: ", err)
	}

	s := service.NewService(&serviceUnit.Service)
	s.Stdout = &logging.Writer{Prefix: serviceName + ": stdout: "}
	s.Stderr = &logging.Writer{Prefix: serviceName + ": stderr: "}

	switch command {
	case start.FullCommand():
		err = s.Start()
		if err != nil {
			log.Fatal(serviceName, ": failed to start service: ", err)
		}
		log.Printf("%s: started with pid %d: %s", serviceName, s.PID(), s.Exe())
	case stop.FullCommand():
		err = s.Stop()
		if err != nil {
			log.Fatal(serviceName, ": failed to stop service: ", err)
		}
	default:
		panic("unexpected command: " + command)
	}
}
