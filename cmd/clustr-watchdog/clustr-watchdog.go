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
	for {
		p.Start()
		p.Wait()
		log.Print(p.Name, ": not running, restarting...")
	}

	panic("unreachable")
}
