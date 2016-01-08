package highservice

import (
	"log"

	"github.com/dmage/clustr/daemon"
	"github.com/dmage/clustr/logging"
	"github.com/dmage/clustr/unit"
)

type HighService struct {
	Name   string
	Unit   *unit.ServiceUnit
	Daemon *daemon.Daemon
}

func HighServiceFromFile(filename string) *HighService {
	serviceName, serviceUnit, err := unit.ServiceUnitFromFile(filename)
	if err != nil {
		log.Fatal("failed to load service: ", err)
	}

	s := daemon.NewDaemon(&serviceUnit.DaemonConfig)
	s.Stdout = &logging.Writer{Prefix: serviceName + ": stdout: "}
	s.Stderr = &logging.Writer{Prefix: serviceName + ": stderr: "}

	err = daemon.LoadState(serviceName, s)
	if err != nil {
		log.Fatal(serviceName, ": failed to load state: ", err)
	}

	return &HighService{
		Name:   serviceName,
		Unit:   serviceUnit,
		Daemon: s,
	}
}

func (hs *HighService) IsRunning() bool {
	alive, err := hs.Daemon.IsRunning()
	if err != nil {
		log.Fatal(hs.Name, ": failed to check service status: ", err)
	}
	return alive
}

func (hs *HighService) Start() {
	if hs.IsRunning() {
		log.Fatal(hs.Name, ": service already running (pid ", hs.Daemon.PID(), ")")
	}

	err := hs.Daemon.Start()
	if err != nil {
		log.Fatal(hs.Name, ": failed to start service: ", err)
	}

	err = daemon.SaveState(hs.Name, hs.Daemon)
	if err != nil {
		log.Fatal(hs.Name, ": failed to save state: ", err)
	}

	log.Printf("%s: started with pid %d: %s", hs.Name, hs.Daemon.PID(), hs.Daemon.Exe())
}

func (hs *HighService) Stop() {
	if !hs.IsRunning() {
		log.Fatal(hs.Name, ": service already stopped")
	}

	log.Printf("%s: stopping pid %d...", hs.Name, hs.Daemon.PID())

	err := hs.Daemon.Stop()
	if err != nil {
		log.Fatal(hs.Name, ": failed to stop service: ", err)
	}
}

func (hs *HighService) Wait() {
	err := hs.Daemon.Wait()
	if err != nil {
		log.Fatal(hs.Name, ": failed to wait service: ", err)
	}
}
