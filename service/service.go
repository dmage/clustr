package service

import (
	"log"

	"github.com/dmage/clustr/daemon"
	"github.com/dmage/clustr/logging"
	"github.com/dmage/clustr/unit"
)

type Service struct {
	Name   string
	Unit   *unit.ServiceUnit
	Daemon *daemon.Daemon
}

func ServiceFromFile(filename string) *Service {
	serviceName, serviceUnit, err := unit.ServiceUnitFromFile(filename)
	if err != nil {
		log.Fatal("failed to load service unit: ", err)
	}

	d := daemon.NewDaemon(&serviceUnit.DaemonConfig)
	d.Stdout = &logging.Writer{Prefix: serviceName + ": stdout: "}
	d.Stderr = &logging.Writer{Prefix: serviceName + ": stderr: "}

	err = daemon.LoadState(serviceName, d)
	if err != nil {
		log.Fatal(serviceName, ": failed to load state: ", err)
	}

	return &Service{
		Name:   serviceName,
		Unit:   serviceUnit,
		Daemon: d,
	}
}

func (s *Service) IsRunning() bool {
	alive, err := s.Daemon.IsRunning()
	if err != nil {
		log.Fatal(s.Name, ": failed to check service status: ", err)
	}
	return alive
}

func (s *Service) Start() {
	if s.IsRunning() {
		log.Fatal(s.Name, ": service already running (pid ", s.Daemon.PID(), ")")
	}

	err := s.Daemon.Start()
	if err != nil {
		log.Fatal(s.Name, ": failed to start service: ", err)
	}

	err = daemon.SaveState(s.Name, s.Daemon)
	if err != nil {
		log.Fatal(s.Name, ": failed to save state: ", err)
	}

	log.Printf("%s: started with pid %d: %s", s.Name, s.Daemon.PID(), s.Daemon.Exe())
}

func (s *Service) Stop() {
	if !s.IsRunning() {
		log.Fatal(s.Name, ": service already stopped")
	}

	log.Printf("%s: stopping pid %d...", s.Name, s.Daemon.PID())

	err := s.Daemon.Stop()
	if err != nil {
		log.Fatal(s.Name, ": failed to stop service: ", err)
	}
}

func (s *Service) Wait() {
	err := s.Daemon.Wait()
	if err != nil {
		log.Fatal(s.Name, ": failed to wait service: ", err)
	}
}
