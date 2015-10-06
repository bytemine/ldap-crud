package slapd

import (
	"os"
	"os/exec"
	"testing"
)

func TestConfig(t *testing.T) {
	config := DefaultConfig
	config.Schemas = []string{"schema/core.schema"}
	_, err := config.Configure()
	if err != nil {
		t.Error(err)
	}

	for _, v := range []string{config.dir, config.db, config.file.Name()} {
		_, err := os.Stat(v)
		if err != nil {
			t.Error(err)
		}
	}

	err = config.Unconfigure()
	if err != nil {
		t.Error(err)
	}

	_, err = config.file.Stat()
	if err == nil {
		t.Error(err)
	}
}

func TestSlapd(t *testing.T) {
	var slapd Slapd

	slapd.Config = &DefaultConfig
	err := slapd.Start()

	if err != nil {
		t.Error(err)
	}

	err = DefaultConfig.Initialize()
	if err != nil {
		t.Error(err)
	}

	// test if we can bind as rootdn
	cmd := exec.Command("ldapwhoami", "-x", "-D", DefaultConfig.Rootdn.Dn, "-w", DefaultConfig.Rootdn.Password, "-H", DefaultConfig.url())
	err = cmd.Run()
	if err != nil {
		t.Error(err)
	}

	err = slapd.Stop()
	if err != nil {
		t.Error(err)
	}

	_, err = slapd.cmd.Process.Wait()
	if err != nil {
		t.Error(err)
	}

}
