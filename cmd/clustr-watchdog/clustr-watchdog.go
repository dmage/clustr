package main

import (
	"log"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/dmage/clustr/watcher"
)

var (
	serviceUnitFile = kingpin.Flag("unit", "Service unit configuration file.").Short('u').Required().String()
)

func main() {
	kingpin.Parse()

	p := watcher.PrisonerFromFile(*serviceUnitFile)
	if p.IsRunning() {
		log.Print(p.Name, ": found process with pid ", p.Service.PID())
	} else {
		p.Start()
	}
	for {
		p.Wait()
		log.Print(p.Name, ": not running, restarting...")
		p.Start()
	}

	panic("unreachable")
}
