package watcher

import (
	"log"

	"github.com/dmage/clustr/logging"
	"github.com/dmage/clustr/service"
	"github.com/dmage/clustr/unit"
)

type Prisoner struct {
	Name    string
	Unit    *unit.ServiceUnit
	Service *service.Service
}

func PrisonerFromFile(filename string) *Prisoner {
	serviceName, serviceUnit, err := unit.ServiceUnitFromFile(filename)
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

	return &Prisoner{
		Name:    serviceName,
		Unit:    serviceUnit,
		Service: s,
	}
}

func (p *Prisoner) IsRunning() bool {
	alive, err := p.Service.IsRunning()
	if err != nil {
		log.Fatal(p.Name, ": failed to check service status: ", err)
	}
	return alive
}

func (p *Prisoner) Start() {
	if p.IsRunning() {
		log.Fatal(p.Name, ": service already running (pid ", p.Service.PID(), ")")
	}

	err := p.Service.Start()
	if err != nil {
		log.Fatal(p.Name, ": failed to start service: ", err)
	}

	err = service.SaveState(p.Name, p.Service)
	if err != nil {
		log.Fatal(p.Name, ": failed to save state: ", err)
	}

	log.Printf("%s: started with pid %d: %s", p.Name, p.Service.PID(), p.Service.Exe())
}

func (p *Prisoner) Stop() {
	if !p.IsRunning() {
		log.Fatal(p.Name, ": service already stopped")
	}

	log.Printf("%s: stopping pid %d...", p.Name, p.Service.PID())

	err := p.Service.Stop()
	if err != nil {
		log.Fatal(p.Name, ": failed to stop service: ", err)
	}
}

func (p *Prisoner) Wait() {
	err := p.Service.Wait()
	if err != nil {
		log.Fatal(p.Name, ": failed to wait service: ", err)
	}
}
