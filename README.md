# ldap-crud
ldap-crud is a collection of Go packages to make handling LDAP easier.

## Packages
### Package crud
This is a simple layer to ease the usage of LDAP with Go. It defines a set of
common operations:
- Create
- Read, ReadAll, ReadAllSiblings, ReadAllSubtree
- Update
- Delete, DeleteSubtree (recursively, not with controls)

### Command schema2go
schema2go generates Go code containing Item definitions usable with package crud.
Note that this is not really polished; ymmv.

### Package slapd
Creates fresh instances of OpenLDAPs slapd for testing purposes.

## Installation
The usual `go get` should work with each of these packages.

