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

	serviceUnit, err := unit.ServiceUnitFromFile("./sleep.service")
	if err != nil {
		panic(err)
	}

	service := service.NewService("sleep.service", &serviceUnit.Service)
	service.WaitDelay = 100 * time.Millisecond
	err = service.Start()
	if err != nil {
		log.Fatal("failed to start service: ", err)
	}

	log.Println("started, sleeping for 2 seconds")
	time.Sleep(2 * time.Second)
	log.Println("stopping")

	err = service.Stop()
	if err != nil {
		log.Fatal("failed to stop service: ", err)
	}

	err = service.Wait()
	if err != nil {
		log.Fatal("failed to wait service: ", err)
	}
}
