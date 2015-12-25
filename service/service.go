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
	Name      string
	WaitDelay time.Duration

	config *Config
	pid    int
	exe    string
}

type Error struct {
	Name string
	Err  error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Name, e.Err)
}

func NewService(name string, config *Config) *Service {
	return &Service{
		Name:   name,
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

	log.Printf("%s: started with pid %d: %s", s.Name, s.pid, s.exe)

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
		return &Error{
			Name: s.Name,
			Err:  fmt.Errorf("stderr pipe: %s", err),
		}
	}
	go LogReader(s.Name+": stderr: ", stderr)

	err = cmd.Run()
	if err != nil {
		return &Error{
			Name: s.Name,
			Err:  err,
		}
	}

	output := stdout.String()
	pid, err := strconv.Atoi(strings.TrimSpace(output))
	if err != nil {
		return &Error{
			Name: s.Name,
			Err:  fmt.Errorf("invalid pid from start command: %s", strconv.Quote(output)),
		}
	}

	return s.initPID(pid)
}

func (s *Service) Stop() error {
	cmd := exec.Command("/bin/sh", "-c", s.config.Stop)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return &Error{
			Name: s.Name,
			Err:  fmt.Errorf("stdout pipe: %s", err),
		}
	}
	go LogReader(s.Name+": stdout: ", stdout)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return &Error{
			Name: s.Name,
			Err:  fmt.Errorf("stderr pipe: %s", err),
		}
	}
	go LogReader(s.Name+": stderr: ", stderr)

	err = cmd.Run()
	if err != nil {
		return &Error{
			Name: s.Name,
			Err:  err,
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
