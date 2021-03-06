package daemon

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"
)

type state struct {
	PID int
	Exe string
}

func SaveState(name string, d *Daemon) error {
	f, err := os.Create(fmt.Sprintf("./.__state__%s.json", name))
	if err != nil {
		return err
	}
	defer f.Close()

	v := state{
		PID: d.PID(),
		Exe: d.Exe(),
	}
	return json.NewEncoder(f).Encode(&v)
}

func LoadState(name string, d *Daemon) error {
	f, err := os.Open(fmt.Sprintf("./.__state__%s.json", name))
	if perr, ok := err.(*os.PathError); ok && perr.Err == syscall.ENOENT {
		return nil
	}
	defer f.Close()

	var v state
	err = json.NewDecoder(f).Decode(&v)
	if err != nil {
		return err
	}

	d.InitPIDExe(v.PID, v.Exe)
	return nil
}
