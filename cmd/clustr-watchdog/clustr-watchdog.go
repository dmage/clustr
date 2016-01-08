package main

import (
	"log"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/dmage/clustr/service"
)

var (
	serviceUnitFile = kingpin.Flag("unit", "Service unit configuration file.").Short('u').Required().String()
)

func main() {
	kingpin.Parse()

	s := service.ServiceFromFile(*serviceUnitFile)
	if s.IsRunning() {
		log.Print(s.Name, ": found process with pid ", s.Daemon.PID())
	} else {
		s.Start()
	}
	for {
		s.Wait()
		log.Print(s.Name, ": not running, restarting...")
		s.Start()
	}

	panic("unreachable")
}
