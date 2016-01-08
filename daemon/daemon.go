package daemon

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Flusher interface {
	Flush() error
}

type WriteFlusher interface {
	io.Writer
	Flusher
}

type Config struct {
	Start string
	Stop  string
}

type Daemon struct {
	WaitDelay time.Duration

	Stdout WriteFlusher
	Stderr WriteFlusher

	config *Config
	pid    int
	exe    string
}

func NewDaemon(config *Config) *Daemon {
	return &Daemon{
		WaitDelay: 100 * time.Millisecond,
		config:    config,
	}
}

func (d *Daemon) initPID(pid int) error {
	d.pid = pid

	exe, err := getProcExe(pid)
	if err != nil {
		return err
	}
	d.exe = exe

	return nil
}

func (d *Daemon) InitPIDExe(pid int, exe string) {
	d.pid = pid
	d.exe = exe
}

func (d *Daemon) IsRunning() (bool, error) {
	if d.pid == 0 || d.exe == "" {
		return false, nil
	}

	err := syscall.Kill(d.pid, 0)
	if err == syscall.ESRCH {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	exe, err := getProcExe(d.pid)
	if err == nil && d.exe != exe {
		return false, nil
	}

	return true, nil
}

func (d *Daemon) Start() error {
	var stdout bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", d.config.Start)
	cmd.Stdout = &stdout
	cmd.Stderr = d.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	if d.Stderr != nil {
		err = d.Stderr.Flush()
		if err != nil {
			return err
		}
	}

	output := stdout.String()
	pid, err := strconv.Atoi(strings.TrimSpace(output))
	if err != nil {
		return fmt.Errorf("invalid pid from start command: %s", strconv.Quote(output))
	}

	return d.initPID(pid)
}

func (d *Daemon) Stop() error {
	cmd := exec.Command("/bin/sh", "-c", d.config.Stop)
	cmd.Stdout = d.Stdout
	cmd.Stderr = d.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	if d.Stdout != nil {
		err = d.Stdout.Flush()
		if err != nil {
			return err
		}
	}

	if d.Stderr != nil {
		err = d.Stderr.Flush()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Daemon) Wait() error {
	for {
		alive, err := d.IsRunning()
		if err != nil {
			return err
		}
		if !alive {
			return nil
		}
		time.Sleep(d.WaitDelay)
	}
}

func (d *Daemon) PID() int {
	return d.pid
}

func (d *Daemon) Exe() string {
	return d.exe
}
