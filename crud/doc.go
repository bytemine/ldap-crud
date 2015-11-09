/*
Package crud implements simple tooling for crud with LDAP.

This is a rather low-level layer to ease the usage of LDAP with Go. It defines a set of
common operations. These operations are performed on objects implementing the "Item"
interface. Commonly, these objects are structs, representing one or more objectclasses
as fields of the struct. The reason for using this approach is to lever Gos type system
to have more safety when handling LDAP. A Item implementation could for example use string
typed fields for SINGLE-VALUE attributes and []string typed fields for multivalued attributes.

An example for an implementation of the Item interface can be
found in the tests.

Note: To run the tests, you must have openldap installed.
*/
package crud
