package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/dmage/clustr/service"

	"github.com/BurntSushi/toml"
	"github.com/samuel/go-zookeeper/zk"
)

type ServiceUnit struct {
	Service service.Config
}

func ServiceUnitFromFile(filename string) (*ServiceUnit, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &ServiceUnit{}
	err = toml.Unmarshal(buf, c)
	return c, err
}

func main() {
	c, _, err := zk.Connect([]string{"clustr_zookeeper_1"}, time.Second)
	if err != nil {
		// panic(err)
	}
	c = c

	serviceUnit, err := ServiceUnitFromFile("./sleep.service")
	if err != nil {
		panic(err)
	}

	service := service.NewService(&serviceUnit.Service)
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
