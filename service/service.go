package service

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func LogReader(prefix string, stream io.Reader) {
	r := bufio.NewReader(stream)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if line != "" {
				log.Print(prefix, line)
			}
			break
		}
		log.Print(prefix, line[0:len(line)-1])
	}
}

type Config struct {
	Start string
	Stop  string
}

type Service struct {
	config *Config
	pid    int
	exe    string
}

func NewService(config *Config) *Service {
	return &Service{
		config: config,
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
	cmd := exec.Command("/bin/sh", "-c", s.config.Start)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %s", err)
	}
	go LogReader(strconv.Quote(s.config.Stop)+": stderr: ", stderr)

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: %s", strconv.Quote(s.config.Start), err)
	}

	output := stdout.String()
	pid, err := strconv.Atoi(strings.TrimSpace(output))
	if err != nil {
		return fmt.Errorf("invalid pid from %s: %s", strconv.Quote(s.config.Start), strconv.Quote(output))
	}

	return s.initPID(pid)
}

func (s *Service) Stop() error {
	cmd := exec.Command("/bin/sh", "-c", s.config.Stop)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe: %s", err)
	}
	go LogReader(strconv.Quote(s.config.Stop)+": stdout: ", stdout)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe: %s", err)
	}
	go LogReader(strconv.Quote(s.config.Stop)+": stderr: ", stderr)

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: %s", strconv.Quote(s.config.Stop), err)
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
		time.Sleep(100 * time.Millisecond)
	}
}
