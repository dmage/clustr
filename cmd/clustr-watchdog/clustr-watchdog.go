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
)

func main() {
	kingpin.Parse()

	serviceName, serviceUnit, err := unit.ServiceUnitFromFile(*serviceUnitFile)
	if err != nil {
		log.Fatal("failed to load service: ", err)
	}

	s := service.NewService(&serviceUnit.Service)
	s.Stdout = &logging.Writer{Prefix: serviceName + ": stdout: "}
	s.Stderr = &logging.Writer{Prefix: serviceName + ": stderr: "}

	err = service.LoadState(serviceName, s)
	if err != nil {
		log.Fatal(serviceName, ": failed to load state: ", err)
	}

	alive, err := s.IsRunning()
	if err != nil {
		log.Fatal(serviceName, ": failed to check service status: ", err)
	}
	if alive {
		log.Fatal(serviceName, ": service already running (pid ", s.PID(), ")")
	}

	for {
		err = s.Start()
		if err != nil {
			log.Fatal(serviceName, ": failed to start service: ", err)
		}
		err = service.SaveState(serviceName, s)
		if err != nil {
			log.Fatal(serviceName, ": failed to save state: ", err)
		}
		log.Printf("%s: started with pid %d: %s", serviceName, s.PID(), s.Exe())

		err = s.Wait()
		if err != nil {
			log.Fatal(serviceName, ": failed to wait service: ", err)
		}
		log.Print(serviceName, ": not running, restarting...")
	}

	panic("unreachable")
}
