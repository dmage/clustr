package unit

import (
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/dmage/clustr/daemon"
)

type ServiceUnit struct {
	DaemonConfig daemon.Config `toml:"Service"`
}

func ServiceUnitFromFile(filename string) (name string, unit *ServiceUnit, err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", nil, err
	}

	c := &ServiceUnit{}
	err = toml.Unmarshal(buf, c)
	return filepath.Base(filename), c, err
}
