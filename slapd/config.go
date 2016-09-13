package slapd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

var DefaultConfig = Config{
	Loglevel: 0,
	Suffix: Object{Dn: "dc=example,dc=com",
		Attributes: map[string][]string{"objectClass": []string{"dcObject", "organization"},
			"o":  []string{"example.com"},
			"dc": []string{"example"}}},
	Rootdn: Object{Dn: "cn=admin,dc=example,dc=com",
		Attributes: map[string][]string{"objectClass": []string{"organizationalRole"},
			"cn": []string{"admin"}},
		Password: "secret"},

	Schemas:        []string{},
	DBType:         "ldif",
	Addr:           "127.0.0.1:9999",
	ConfigTemplate: DefaultConfigTemplate,
	LdifTemplate:   DefaultLdifTemplate,
}

// Default LDIF generation template, used to add the suffix and root objects with ldapadd
var DefaultLdifTemplate = `
{{range $i, $o := . }}dn: {{ $o.Dn }}
{{range $k, $vs := $o.Attributes}}{{range $i, $v := $vs }}{{ $k }}: {{ $v }}
{{end}}{{end}}
{{end}}

`

// Default slapd config template
var DefaultConfigTemplate = `
{{range $i, $v := .Schemas}}include {{$v}}
{{end}}

database {{.DBType}}
suffix {{.Suffix}}
rootdn {{.Rootdn}}
rootpw {{.Rootpw}}
directory {{.Db}}

access to attrs=userPassword
        by self write
        by anonymous auth
        by users none

access to * by * read
`

// Simple representation of an LDAP object
type Object struct {
	Dn         string
	Attributes map[string][]string
	Password   string
}

// Default Configurer implementation for a slapd instance
type Config struct {
	// Loglevel (be aware that high loglevels generate really much output, which may clog the pipes).
	Loglevel int

	// listen address and port of slapd in the format "host:port"
	Addr string

	// base dn suffix
	Suffix Object

	//root object dn
	Rootdn Object

	//schema files to include
	Schemas []string

	// database type
	DBType string

	// slapd.conf template
	ConfigTemplate string

	// LDIF template to use for adding suffix and root objects. Used with ldapadd.
	LdifTemplate string

	// base dir for slapd config and db
	dir string

	// database directory
	db string

	// slapd config file
	file *os.File
}

// url returns an url to connect to the slapd in the format "ldap://address"
func (c *Config) url() string {
	return fmt.Sprintf("ldap://%v", c.Addr)
}

// Create a configuration and a slapd process struct which uses this config
func (c *Config) Configure() (*exec.Cmd, error) {
	var err error
	c.dir, err = ioutil.TempDir(os.TempDir(), "slapd")
	if err != nil {
		return nil, err
	}

	c.db = filepath.Join(c.dir, "db")
	err = os.Mkdir(c.db, os.ModeDir|0777)
	if err != nil {
		return nil, err
	}

	c.file, err = ioutil.TempFile(c.dir, "slapd")

	t := template.Must(template.New("slapdconfig").Parse(c.ConfigTemplate))

	templateConfig := struct {
		Schemas []string
		DBType  string
		Suffix  string
		Rootdn  string
		Rootpw  string
		Db      string
	}{Schemas: c.Schemas, DBType: c.DBType, Suffix: c.Suffix.Dn, Rootdn: c.Rootdn.Dn, Rootpw: c.Rootdn.Password, Db: c.db}

	err = t.Execute(c.file, templateConfig)
	if err != nil {
		return nil, err
	}

	err = c.file.Close()
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("slapd", "-d", fmt.Sprintf("%v", c.Loglevel), "-h", c.url(), "-f", c.file.Name())

	return cmd, nil
}

// Unconfigure deletes the slapd directory created by Configure
func (c *Config) Unconfigure() error {
	return os.RemoveAll(c.dir)
}

// Initialize adds entries for the base and the root object
func (c *Config) Initialize() error {
	objects := make([]Object, 2)
	objects[0] = c.Suffix
	objects[1] = c.Rootdn

	t := template.Must(template.New("ldif").Parse(c.LdifTemplate))

	cmd := exec.Command("ldapadd", "-x", "-D", c.Rootdn.Dn, "-w", c.Rootdn.Password, "-H", c.url())

	cmd.Stderr = os.Stdout

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = t.Execute(stdinPipe, objects)
	if err != nil {
		return err
	}

	stdinPipe.Close()

	return cmd.Wait()
}

// Clean removes all entries from the ldap. You have to run Initialize() again to re-add the admin entry.
func (c *Config) Clean() error {
	cmd := exec.Command("ldapdelete", "-r", "-x", "-D", c.Rootdn.Dn, "-w", c.Rootdn.Password, "-H", c.url(), c.Suffix.Dn)
	return cmd.Run()
}
