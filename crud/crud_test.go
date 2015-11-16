package crud

import (
	"errors"
	"fmt"
	"github.com/bytemine/ldap-crud/slapd"
	"github.com/rbns/ldap"
	"testing"
)

var foobarPerson = Person{sn: []string{"Foobar"}}
var fritzFoobarPerson = Person{sn: []string{"Foobar"}, cn: []string{"Fritz"}}
var fritzBarbazPerson = Person{sn: []string{"Bazbar"}, cn: []string{"Fritz"}}
var fritzQuxPerson = Person{sn: []string{"Qux"}, cn: []string{"Fritz"}, dn: "sn=Qux,sn=Foobar"}
var gonzoPerson = Person{sn: []string{"Foobar"}, cn: []string{"Gonzo", "von"}}

// Standard LDAP person with its must attributes
type Person struct {
	dn string
	sn []string
	cn []string
}

func NewPerson(sn, cn []string) *Person {
	return &Person{sn: sn, cn: cn}
}

// Return a copy of itself
func (p *Person) Copy() Item {
	return &Person{sn: p.sn, cn: p.cn}
}

// We've decided that "sn" should be the attribute used in the DN.
// In a more real-world application the Person struct would
// have a member of type Item which represents the parent and prepending
// it's own rdn to the dn of the parent.
func (p *Person) Dn() string {
	if p.dn == "" {
		return fmt.Sprintf("sn=%v", p.sn[0])
	}
	return p.dn
}

func (p *Person) FilterObjectClass() string {
	return "person"
}

func (p *Person) MarshalLDAP() (*ldap.Entry, error) {
	if len(p.sn) < 1 {
		return nil, errors.New("sn is a must attribute")
	}
	if len(p.cn) < 1 {
		return nil, errors.New("cn is a must attribute")
	}

	entry := ldap.NewEntry("")
	entry.AddAttributeValue("objectClass", "person")

	entry.AddAttributeValues("sn", p.sn)
	entry.AddAttributeValues("cn", p.cn)
	entry.AddAttributeValue("userPassword", "foobar")

	return entry, nil
}

func (p *Person) UnmarshalLDAP(entry *ldap.Entry) error {
	p.dn = entry.DN
	p.sn = entry.GetAttributeValues("sn")
	if len(p.sn) < 1 {
		return errors.New("sn is a must attribute")
	}

	p.cn = entry.GetAttributeValues("cn")
	if len(p.cn) < 1 {
		return errors.New("cn is a must attribute")
	}

	return nil
}

func TestAppendBaseDn(t *testing.T) {
	c := New(nil, "dc=example,dc=com")

	if c.appendBaseDn("cn=foo") != "cn=foo,dc=example,dc=com" {
		t.Fail()
	}
}

func TestRemoveBaseDn(t *testing.T) {
	c := New(nil, "dc=example,dc=com")

	if c.removeBaseDn("cn=foo,dc=example,dc=com") != "cn=foo" {
		t.Fail()
	}
}

func TestParentDn(t *testing.T) {
	if parentDn("cn=foo,dc=example,dc=com") != "dc=example,dc=com" {
		t.Fail()
	}
}

func testCreate(t *testing.T) {
	lc := ldap.NewConnection(slapd.DefaultConfig.Address())
	err := lc.Connect()
	if err != nil {
		t.Error(err)
	}

	err = lc.Bind(slapd.DefaultConfig.Rootdn.Dn, slapd.DefaultConfig.Rootdn.Password)

	if err != nil {
		t.Error(err)
	}

	c := New(lc, "dc=example,dc=com")

	err = c.Create(&fritzFoobarPerson)
	t.Logf("Created: %+v", fritzFoobarPerson)

	if err != nil {
		t.Error(err)
	}

	p := Person{sn: []string{"Foobar"}, cn: []string{}}
	err = c.Read(&p)
	if err != nil {
		t.Error("object wasn't deleted from ldap")
	}

	t.Logf("Read: %+v", p)
}

func TestCreate(t *testing.T) {
	var s = new(slapd.Slapd)

	s.Config = &slapd.DefaultConfig
	err := s.StartAndInitialize()
	defer s.Stop()
	if err != nil {
		t.Error(err)
	}

	testCreate(t)
}

func testRead(t *testing.T) {
	lc := ldap.NewConnection("localhost:9999")
	err := lc.Connect()
	if err != nil {
		t.Error(err)
	}

	err = lc.Bind(slapd.DefaultConfig.Rootdn.Dn, slapd.DefaultConfig.Rootdn.Password)
	if err != nil {
		t.Error(err)
	}

	c := New(lc, "dc=example,dc=com")

	p := foobarPerson

	err = c.Read(&p)
	t.Logf("Read: %+v", p)
	if err != nil {
		t.Error(err)
	}
}

func TestRead(t *testing.T) {
	var s = new(slapd.Slapd)

	s.Config = &slapd.DefaultConfig
	err := s.StartAndInitialize()
	defer s.Stop()
	if err != nil {
		t.Error(err)
	}

	testCreate(t)
	testRead(t)
}

func testUpdate(t *testing.T) {
	lc := ldap.NewConnection("localhost:9999")
	err := lc.Connect()
	if err != nil {
		t.Error(err)
	}

	err = lc.Bind(slapd.DefaultConfig.Rootdn.Dn, slapd.DefaultConfig.Rootdn.Password)
	if err != nil {
		t.Error(err)
	}

	c := New(lc, "dc=example,dc=com")

	// get the person created in testRead
	oldPerson := foobarPerson

	err = c.Read(&oldPerson)
	if err != nil {
		t.Error(err)
	}

	t.Logf("values before update: %+v", oldPerson)

	// create a new person with new cn values "Gonzo", "von"
	// and perform an update
	newPerson := gonzoPerson

	err = c.Update(&newPerson)
	if err != nil {
		t.Error(err)
	}

	// get the person again from ldap. it should now have the cn values
	// we previously set
	err = c.Read(&oldPerson)
	if err != nil {
		t.Error(err)
	}

	t.Logf("values after update: %+v", oldPerson)

	if !equalStringSlice(oldPerson.sn, gonzoPerson.sn) || !equalStringSlice(oldPerson.cn, gonzoPerson.cn) {
		t.Error("Read entry has unexpected attribute values. Expected:", gonzoPerson, "Got:", oldPerson)
	}
}

func TestUpdate(t *testing.T) {
	/*	var s = new(slapd.Slapd)

		s.Config = &slapd.DefaultConfig */
	s := slapd.New(nil)
	err := s.StartAndInitialize()
	defer s.Stop()

	if err != nil {
		t.Error(err)
	}

	testCreate(t)
	testUpdate(t)
}

func testReadAll(t *testing.T) {
	lc := ldap.NewConnection("localhost:9999")
	err := lc.Connect()
	if err != nil {
		t.Error(err)
	}

	err = lc.Bind(slapd.DefaultConfig.Rootdn.Dn, slapd.DefaultConfig.Rootdn.Password)
	if err != nil {
		t.Error(err)
	}

	c := New(lc, "dc=example,dc=com")

	err = c.Create(&fritzFoobarPerson)
	if err != nil {
		t.Error(err)
	}

	err = c.Create(&fritzBarbazPerson)
	if err != nil {
		t.Error(err)
	}

	entries, err := c.ReadAll(&foobarPerson, "", ScopeWholeSubtree, "(objectClass=%v)", "person")
	if err != nil {
		t.Error(err)
	}

	if len(entries) != 2 {
		t.Error("Expected exactly two results, got", len(entries))
	}

	var count int
	for _, v := range entries {
		t.Log(v.Dn(), v)
		for _, w := range []Person{fritzFoobarPerson, fritzBarbazPerson} {
			if equalStringSlice(v.(*Person).sn, w.sn) && equalStringSlice(v.(*Person).cn, w.cn) {
				count++
				continue
			}
		}
	}

	if count != 2 {
		t.Fail()
	}
}

func TestReadAll(t *testing.T) {
	var s = new(slapd.Slapd)

	s.Config = &slapd.DefaultConfig
	err := s.StartAndInitialize()
	if err != nil {
		t.Error(err)
	}
	defer s.Stop()

	testReadAll(t)
}

func testReadAllSiblings(t *testing.T) {
	lc := ldap.NewConnection("localhost:9999")
	err := lc.Connect()
	if err != nil {
		t.Error(err)
	}

	err = lc.Bind(slapd.DefaultConfig.Rootdn.Dn, slapd.DefaultConfig.Rootdn.Password)
	if err != nil {
		t.Error(err)
	}

	c := New(lc, "dc=example,dc=com")

	err = c.Create(&fritzFoobarPerson)
	if err != nil {
		t.Error(err)
	}

	err = c.Create(&fritzBarbazPerson)
	if err != nil {
		t.Error(err)
	}

	entries, err := c.ReadAllSiblings(&foobarPerson)
	if err != nil {
		t.Error(err)
	}

	if len(entries) != 2 {
		t.Error("Expected exactly two results, got", len(entries))
	}

	var count int
	for _, v := range entries {
		t.Log(v.Dn(), v)
		for _, w := range []Person{fritzFoobarPerson, fritzBarbazPerson} {
			if equalStringSlice(v.(*Person).sn, w.sn) && equalStringSlice(v.(*Person).cn, w.cn) {
				count++
				continue
			}
		}
	}

	if count != 2 {
		t.Fail()
	}
}

func TestReadAllSiblings(t *testing.T) {
	var s = new(slapd.Slapd)

	s.Config = &slapd.DefaultConfig
	err := s.StartAndInitialize()
	if err != nil {
		t.Error(err)
	}
	defer s.Stop()

	testReadAllSiblings(t)
}

func testReadAllSubtree(t *testing.T) {
	lc := ldap.NewConnection("localhost:9999")
	err := lc.Connect()
	if err != nil {
		t.Error(err)
	}

	err = lc.Bind(slapd.DefaultConfig.Rootdn.Dn, slapd.DefaultConfig.Rootdn.Password)
	if err != nil {
		t.Error(err)
	}

	c := New(lc, "dc=example,dc=com")

	err = c.Create(&fritzFoobarPerson)
	if err != nil {
		t.Error(err)
	}

	err = c.Create(&fritzQuxPerson)
	if err != nil {
		t.Error(err)
	}

	entries, err := c.ReadAllSubtree(&foobarPerson)
	if err != nil {
		t.Error(err)
	}

	if len(entries) != 2 {
		t.Error("Expected exactly two results, got", len(entries))
	}

	var count int
	for _, v := range entries {
		t.Log(v.Dn(), v)
		for _, w := range []*Person{&fritzFoobarPerson, &fritzQuxPerson} {
			if equalStringSlice(v.(*Person).sn, w.sn) && equalStringSlice(v.(*Person).cn, w.cn) {
				count++
				continue
			}
		}
	}

	if count != 2 {
		t.Fail()
	}
}

func TestReadAllSubtree(t *testing.T) {
	var s = new(slapd.Slapd)

	s.Config = &slapd.DefaultConfig
	err := s.StartAndInitialize()
	if err != nil {
		t.Error(err)
	}
	defer s.Stop()

	testReadAllSubtree(t)
}
func testDelete(t *testing.T) {
	lc := ldap.NewConnection("localhost:9999")
	err := lc.Connect()
	if err != nil {
		t.Error(err)
	}

	err = lc.Bind(slapd.DefaultConfig.Rootdn.Dn, slapd.DefaultConfig.Rootdn.Password)
	if err != nil {
		t.Error(err)
	}

	c := New(lc, "dc=example,dc=com")

	p := foobarPerson

	err = c.Delete(&p)
	if err != nil {
		t.Error(err)
	}

	err = c.Read(&p)
	if err == nil {
		t.Error("object wasn't deleted from ldap")
	}
}

func TestDelete(t *testing.T) {
	var s = new(slapd.Slapd)

	s.Config = &slapd.DefaultConfig
	err := s.StartAndInitialize()
	defer s.Stop()
	if err != nil {
		t.Error(err)
	}

	testCreate(t)
	testDelete(t)
}

func TestDeleteSubtree(t *testing.T) {
	var s = new(slapd.Slapd)
	s.Config = &slapd.DefaultConfig
	err := s.StartAndInitialize()
	defer s.Stop()
	if err != nil {
		t.Error(err)
	}

	lc := ldap.NewConnection("localhost:9999")
	err = lc.Connect()
	if err != nil {
		t.Error(err)
	}

	err = lc.Bind(slapd.DefaultConfig.Rootdn.Dn, slapd.DefaultConfig.Rootdn.Password)
	if err != nil {
		t.Error(err)
	}

	c := New(lc, "dc=example,dc=com")

	err = c.Create(&fritzFoobarPerson)
	if err != nil {
		t.Error(err)
	}

	var fritzSubPerson = fritzFoobarPerson
	fritzSubPerson.dn = fmt.Sprintf("sn=%v,%v", fritzSubPerson.sn, fritzFoobarPerson.Dn())

	err = c.Create(&fritzSubPerson)
	if err != nil {
		t.Error(err)
	}

	err = c.DeleteSubtree(&fritzFoobarPerson)
	if err != nil {
		t.Error(err)
	}

	err = c.Read(&fritzSubPerson)
	if err == nil {
		t.Error("object wasn't deleted from ldap")
	}

	err = c.Read(&fritzFoobarPerson)
	if err == nil {
		t.Error("object wasn't deleted from ldap")
	}
}

func TestPasswd(t *testing.T) {
	var s = new(slapd.Slapd)
	s.Config = &slapd.DefaultConfig
	err := s.StartAndInitialize()
	defer s.Stop()
	if err != nil {
		t.Error(err)
	}

	lc := ldap.NewConnection("localhost:9999")
	err = lc.Connect()
	if err != nil {
		t.Error(err)
	}

	err = lc.Bind(slapd.DefaultConfig.Rootdn.Dn, slapd.DefaultConfig.Rootdn.Password)
	if err != nil {
		t.Error(err)
	}

	c := New(lc, "dc=example,dc=com")

	// create test person
	err = c.Create(&fritzFoobarPerson)
	if err != nil {
		t.Error(err)
	}

	// set password of test person to "foobaz"
	err = c.Passwd(&fritzFoobarPerson, "foobaz")
	if err != nil {
		t.Error(err)
	}

	c.Close()

	lc = ldap.NewConnection("localhost:9999")
	err = lc.Connect()
	if err != nil {
		t.Error(err)
	}

	// try to login as the test person with password "foobaz"
	err = lc.Bind(fritzFoobarPerson.Dn()+","+slapd.DefaultConfig.Suffix.Dn, "foobaz")
	if err != nil {
		t.Error(err)
	}

	c = New(lc, "dc=example,dc=com")

	// let the test person change its own password to "foobar"
	// this needs these acls set in slapd.conf:
	// access to attrs=userPassword
	//	by self write
	//	by anonymous auth
	//	by users none
	// access to * by * read

	err = c.Passwd(nil, "foobar")
	if err != nil {
		t.Error(err)
	}

	c.Close()

	lc = ldap.NewConnection("localhost:9999")
	err = lc.Connect()
	if err != nil {
		t.Error(err)
	}

	// try to login as the test person with password "foobar"
	err = lc.Bind(fritzFoobarPerson.Dn()+","+slapd.DefaultConfig.Suffix.Dn, "foobar")
	if err != nil {
		t.Error(err)
	}

}
