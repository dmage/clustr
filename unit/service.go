package unit

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"

	"github.com/dmage/clustr/service"
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
