package slapd

import (
	"bytes"
	"net"
	"os"
	"path"
	"os/exec"
	"time"
)

// Export our schema directory.
var Schemadir = path.Join(os.Getenv("GOPATH"), "src/github.com/bytemine/ldap-crud/slapd/schema")

// Slapd represents a slapd instance and its configuration.
type Slapd struct {
	// Configuration of the slapd instance
	Config Configurer

	// slapd stdout and stderr are written to this buffer
	Output bytes.Buffer

	// slapd process
	cmd *exec.Cmd
}

// New creates a new Slapd using the specified configuration.
// If config is nil, the DefaultConfig is used as configuration.
func New(config Configurer) *Slapd {
	s := new(Slapd)

	if config == nil {
		config = &DefaultConfig
	}

	s.Config = config
	return s
}

// Configure and start a slapd process
func (s *Slapd) Start() error {
	cmd, err := s.Config.Configure()
	if err != nil {
		return err
	}
	s.cmd = cmd
	s.cmd.Stdout = &s.Output
	s.cmd.Stderr = &s.Output

	err = s.cmd.Start()
	if err != nil {
		return err
	}

	// Wait for slapd to start.
	err = s.wait()

	return err
}

// wait waits for slapd to become ready.
//
// The waiting is done by a rather "brute-force" method:
// Repeatedly try to connect to the slapd process, abort after a fixed maximum number of
// tries.
func (s *Slapd) wait() error {
	var err error
	for i := uint(0); i < s.Config.Maxtries(); i++ {
		conn, err := net.Dial("tcp", s.Config.Address())
		if err == nil {
			// make sure it really has enough time
			time.Sleep(50 * time.Millisecond)
			defer conn.Close()
			break
		}
		// ease the load
		time.Sleep(10 * time.Millisecond)
	}

	return err
}

// The same as slapd.Start but also runs the Initialize method of the associated config
func (s *Slapd) StartAndInitialize() error {
	err := s.Start()
	if err != nil {
		return err
	}

	err = s.Config.Initialize()
	if err != nil {
		return err
	}
	return nil
}

// Stop the slapd process and clean up
func (s *Slapd) Stop() error {
	err := s.cmd.Process.Kill()
	if err != nil {
		return err
	}

	// wait a bit, in case a slapd shuld be started again short after stopping this instance to that the
	// network port is available etc.
	time.Sleep(50 * time.Millisecond)

	err = s.Config.Unconfigure()

	return err
}
