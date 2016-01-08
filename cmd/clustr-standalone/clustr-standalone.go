package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/dmage/clustr/service"
)

var (
	servicesDir = kingpin.Flag("services", "Path to directory with services.").Short('s').Required().String()
)

func main() {
	kingpin.Parse()

	files, err := ioutil.ReadDir(*servicesDir)
	if err != nil {
		log.Fatalf("failed to read directory %s: %s", *servicesDir, err)
	}
	for _, f := range files {
		filename := f.Name()

		if !strings.HasSuffix(filename, ".service") {
			continue
		}

		go func(serviceUnitFile string) {
			s := service.ServiceFromFile(serviceUnitFile)
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
		}(filepath.Join(*servicesDir, filename))
	}

	wait := make(chan struct{})
	<-wait
}
