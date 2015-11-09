package crud

import (
	"errors"
	"fmt"
	"github.com/rbns/ldap"
	"log"
	"strings"
)

// Scope is a clone of the Scope constants of the ldap package for using with ReadAll.
type Scope ldap.Scope

const (
	ScopeBaseObject   Scope = 0
	ScopeSingleLevel  Scope = 1
	ScopeWholeSubtree Scope = 2
)

// Item is the interface implemented by objects that can be used for CRUD
type Item interface {
	// Returns a new Item of the same type and contents.
	// The problems this method solve can maybe also solved by using
	// reflect, but I don't want to open that can of worms.
	Copy() Item

	// Method to marshal itself to an *ldap.Entry
	MarshalLDAP() (*ldap.Entry, error)

	// Method to unmarshal an *ldap.Entry into itself
	UnmarshalLDAP(*ldap.Entry) error

	// Returns the DN of the object
	Dn() string

	// Returns the object class which should be used when searching for all objects of this kind
	FilterObjectClass() string
}

// A Manager performs the CRUD operations for objects implementing Item
type Manager struct {
	// Debug mode flag
	Debug bool

	// base DN to append
	baseDn string

	// Connection to use
	conn *ldap.Connection
}

// New creates a new Manager.
//
// The supplied Connection has to be connected and if necessary
// bound.
func New(c *ldap.Connection, baseDn string) *Manager {
	return &Manager{Debug: false, conn: c, baseDn: baseDn}
}

// Close closes a Manger and its connections, preventing further usage.
func (c *Manager) Close() error {
	return c.conn.Close()
}

// Appends the baseDn if it is not empty, otherwise the
// dn is returned unmodified
func (c *Manager) appendBaseDn(dn string) string {
	// baseDn is not set
	if c.baseDn == "" {
		return dn
	}

	// the supplied dn is empty, return the baseDn without leading ","
	if dn == "" {
		return c.baseDn
	}

	// return concatenated dn and baseDn
	return fmt.Sprintf("%v,%v", dn, c.baseDn)
}

// Removes the baseDn if it is not empty, otherwise the
// dn is returned unmodified
func (c *Manager) removeBaseDn(dn string) string {
	if c.baseDn == "" {
		return dn
	}

	return dn[0 : len(dn)-len(c.baseDn)-1]
}

// Create item in LDAP
func (c *Manager) Create(item Item) error {
	entry, err := item.MarshalLDAP()
	if err != nil {
		return err
	}

	addRequest := ldap.NewAddRequest(c.appendBaseDn(item.Dn()))
	addRequest.Entry.Attributes = entry.Attributes

	if c.Debug {
		log.Println("Add request:", addRequest)
	}

	return c.conn.Add(addRequest)
}

// Read values for the attributes of item from LDAP
func (c *Manager) Read(item Item) error {
	searchRequest := ldap.NewSimpleSearchRequest(c.appendBaseDn(item.Dn()), ldap.ScopeBaseObject, "(objectClass=*)", nil)

	if c.Debug {
		log.Println("Search request:", searchRequest)
	}

	results, err := c.conn.Search(searchRequest)
	if err != nil {
		return err
	}

	if len(results.Entries) == 0 {
		return errors.New("No search results.")
	} else if len(results.Entries) > 1 {
		return errors.New("More than one search result.")
	}

	results.Entries[0].DN = c.removeBaseDn(results.Entries[0].DN)

	return item.UnmarshalLDAP(results.Entries[0])
}

// ReadAll searches for all objects which are of the same type as item and match the criteria.
// dn is the root of the subtree which is searched, scope is the scope of the search, filter
// is an ldap filter string.
func (c *Manager) ReadAll(item Item, dn string, scope Scope, filter string) ([]Item, error) {
	searchRequest := ldap.NewSimpleSearchRequest(c.appendBaseDn(dn), ldap.Scope(scope), filter, nil)

	if c.Debug {
		log.Println("Search Request:", searchRequest)
	}

	results, err := c.conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	items := make([]Item, len(results.Entries))
	for i, v := range results.Entries {
		v.DN = c.removeBaseDn(v.DN)
		items[i] = item.Copy()
		items[i].UnmarshalLDAP(v)
		if err != nil {
			return nil, err
		}
	}

	return items, nil
}

// parentDn returns the dn of the parent item. it does so by removing the first
// of comma-seperated fields of dn. The resulting dn may be the empty string.
func parentDn(dn string) string {
	return strings.Join(strings.Split(dn, ",")[1:], ",")
}

// ReadAllSiblings searches for all objects which are:
// a) Of the same kind as item
// b) On the same level as item
func (c *Manager) ReadAllSiblings(item Item) ([]Item, error) {
	filter := fmt.Sprintf("(objectClass=%v)", item.FilterObjectClass())
	return c.ReadAll(item, parentDn(item.Dn()), ScopeSingleLevel, filter)
}

// ReadAllSubtree searches for all objects which are:
// a) Of the same kinde as item
// b) On the same level and below item
func (c *Manager) ReadAllSubtree(item Item) ([]Item, error) {
	filter := fmt.Sprintf("(objectClass=%v)", item.FilterObjectClass())
	return c.ReadAll(item, parentDn(item.Dn()), ScopeWholeSubtree, filter)
}

// Are two string slices equal?
func equalStringSlice(a, b []string) bool {
	if len(a) == len(b) {
		for i, v := range a {
			if b[i] != v {
				return false
			}
		}
		return true
	}

	return false
}

// Build list of ldap modification operations
func (c *Manager) newModifyRequest(oldItem Item, newItem Item) (*ldap.ModifyRequest, error) {
	modifyRequest := ldap.NewModifyRequest(c.appendBaseDn(oldItem.Dn()))

	oldEntry, err := oldItem.MarshalLDAP()
	if err != nil {
		return nil, err
	}

	newEntry, err := newItem.MarshalLDAP()
	if err != nil {
		return nil, err
	}

	// add all new or modified attributes to modifyRequest
	for _, v := range newEntry.Attributes {
		oldValues := oldEntry.GetAttributeValues(v.Name)
		newValues := newEntry.GetAttributeValues(v.Name)

		// if the attribute didn't exist previously, create it
		if len(oldValues) == 0 {
			modifyRequest.AddMod(ldap.NewMod(ldap.ModAdd, v.Name, newValues))
		} else {
			// if the attribute existed, add a modification operation if the new values differ
			if !equalStringSlice(oldValues, newValues) {
				modifyRequest.AddMod(ldap.NewMod(ldap.ModReplace, v.Name, newValues))
			}
		}

	}

	// add all removed attributes to modifyRequest
	for _, v := range oldEntry.Attributes {
		oldValues := oldEntry.GetAttributeValues(v.Name)
		newValues := newEntry.GetAttributeValues(v.Name)

		// if an attribute existed in oldEntry but not in newEntry, delete it
		if len(newValues) == 0 {
			modifyRequest.AddMod(ldap.NewMod(ldap.ModDelete, v.Name, oldValues))
		}
	}

	return modifyRequest, nil
}

// Update the LDAP attributes of Item
func (c *Manager) Update(newItem Item) error {

	// get the values currently stored in ldap
	oldItem := newItem.Copy()

	err := c.Read(oldItem)
	if err != nil {
		return err
	}

	if c.Debug {
		log.Printf("Old item: %+v\n", oldItem)
		log.Printf("New item: %+v\n", newItem)
	}

	modifyRequest, err := c.newModifyRequest(oldItem, newItem)
	if err != nil {
		return err
	}

	if c.Debug {
		log.Println("Modify request:", modifyRequest)
	}

	return c.conn.Modify(modifyRequest)
}

// Delete an item
func (c *Manager) Delete(item Item) error {
	deleteRequest := ldap.NewDeleteRequest(c.appendBaseDn(item.Dn()))

	if c.Debug {
		log.Println("Delete Request:", deleteRequest)
	}

	return c.conn.Delete(deleteRequest)
}

// Helper method to recursively delete a subtree
func (c *Manager) deleteRecursive(dn string) error {
	// first recursively delete all subentrys
	searchRequest := ldap.NewSimpleSearchRequest(dn, ldap.ScopeSingleLevel, "(objectClass=*)", nil)

	if c.Debug {
		log.Println("Search Request:", searchRequest)
	}

	results, err := c.conn.Search(searchRequest)
	if err != nil {
		return err
	}

	for _, v := range results.Entries {
		err = c.deleteRecursive(v.DN)
		if err != nil {
			return err
		}
	}

	// delete the root of the current tree
	deleteRequest := ldap.NewDeleteRequest(dn)
	return c.conn.Delete(deleteRequest)
}

// DeleteSubtree recursively deletes a subtree without using special controls.
func (c *Manager) DeleteSubtree(item Item) error {
	return c.deleteRecursive(c.appendBaseDn(item.Dn()))
}

// Passwd changes the password of a dn.
func (c *Manager) Passwd(item Item, passwd string) error {
	var dn string
	if item == nil {
		dn = ""
	} else {
		dn = c.appendBaseDn(item.Dn())
	}
	passwdRequest := ldap.PasswordModifyRequest{dn, "", passwd}

	if c.Debug {
		log.Println("Password Request:", passwdRequest)
	}

	return c.conn.Passwd(&passwdRequest)
}
