package service

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

type Service struct {
	WaitDelay time.Duration

	Stdout WriteFlusher
	Stderr WriteFlusher

	config *Config
	pid    int
	exe    string
}

func NewService(config *Config) *Service {
	return &Service{
		WaitDelay: 100 * time.Millisecond,
		config:    config,
	}
}

func (s *Service) initPID(pid int) error {
	s.pid = pid

	exe, err := getProcExe(pid)
	if err != nil {
		return err
	}
	s.exe = exe

	return nil
}

func (s *Service) InitPIDExe(pid int, exe string) {
	s.pid = pid
	s.exe = exe
}

func (s *Service) IsRunning() (bool, error) {
	if s.pid == 0 || s.exe == "" {
		return false, nil
	}

	err := syscall.Kill(s.pid, 0)
	if err == syscall.ESRCH {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	exe, err := getProcExe(s.pid)
	if err == nil && s.exe != exe {
		return false, nil
	}

	return true, nil
}

func (s *Service) Start() error {
	var stdout bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", s.config.Start)
	cmd.Stdout = &stdout
	cmd.Stderr = s.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	if s.Stderr != nil {
		err = s.Stderr.Flush()
		if err != nil {
			return err
		}
	}

	output := stdout.String()
	pid, err := strconv.Atoi(strings.TrimSpace(output))
	if err != nil {
		return fmt.Errorf("invalid pid from start command: %s", strconv.Quote(output))
	}

	return s.initPID(pid)
}

func (s *Service) Stop() error {
	cmd := exec.Command("/bin/sh", "-c", s.config.Stop)
	cmd.Stdout = s.Stdout
	cmd.Stderr = s.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	if s.Stdout != nil {
		err = s.Stdout.Flush()
		if err != nil {
			return err
		}
	}

	if s.Stderr != nil {
		err = s.Stderr.Flush()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Wait() error {
	for {
		alive, err := s.IsRunning()
		if err != nil {
			return err
		}
		if !alive {
			return nil
		}
		time.Sleep(s.WaitDelay)
	}
}

func (s *Service) PID() int {
	return s.pid
}

func (s *Service) Exe() string {
	return s.exe
}
