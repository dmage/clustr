package main

import (
	"log"
	"time"

	"github.com/samuel/go-zookeeper/zk"

	"github.com/dmage/clustr/service"
	"github.com/dmage/clustr/unit"
)

func main() {
	c, _, err := zk.Connect([]string{"clustr_zookeeper_1"}, time.Second)
	if err != nil {
		// panic(err)
	}
	c = c

	serviceName, serviceUnit, err := unit.ServiceUnitFromFile("./sleep.service")
	if err != nil {
		panic(err)
	}

	s := service.NewService(&serviceUnit.Service)
	err = s.Start()
	if err != nil {
		log.Fatal(serviceName, ": failed to start service: ", err)
	}

	log.Printf("started with pid %d: %s", s.PID(), s.Exe())
	log.Println("sleeping for 2 seconds")
	time.Sleep(2 * time.Second)
	log.Println("stopping")

	err = s.Stop()
	if err != nil {
		log.Fatal(serviceName, ": failed to stop service: ", err)
	}

	err = s.Wait()
	if err != nil {
		log.Fatal(serviceName, ": failed to wait service: ", err)
	}
}
