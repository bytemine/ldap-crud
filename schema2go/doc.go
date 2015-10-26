/*
Command schema2go generates Go-code from object definitions and ldap schema files.
The resulting code consists of structs which fulfill the Item interface of package crud [1].

The code is either written to a file or, when outfile is the empty
string "", to stdout.

An example object definitions file could look like this:
	[
	{		"Name": "Account",
			"Desc": "a simple account",
			"ObjectClasses": ["posixAccount"],
			"FilterObjectClass": "posixAccount"
	}
	]
Where the "Name" field is the name of the struct to be generated. The name MUST be upper case
as it would otherwise not be exported. "Desc" is a description which is included as documentation
comment for this struct. "ObjectClasses" is a list of object classes defined in the schema files
this struct should include. "FilterObjectClass" is the "main" object class that should be used
for searching by the crud package.

If this were the content of objects.json, a call
	schema2go -pkg account -obj objects.json -out account.go nis.schema
would write Go-code for package company containing an "Account" struct and methods to marshal
and unmarshal it from an ldap.Entry [2][3] to account.go .

For the above example the resulting struct definition would look like this (code and comments omitted):
	type Account struct {
		dn           string
		PosixAccount bool

		//MUST attributes
		Cn []string `json:"omitempty"`
		GidNumber string `json:"omitempty"`
		HomeDirectory string `json:"omitempty"`
		Uid []string `json:"omitempty"`
		UidNumber string `json:"omitempty"`

		// MAY attributes
		Description []string `json:"omitempty"`
		Gecos string `json:"omitempty"`
		LoginShell string `json:"omitempty"`
		UserPassword []string `json:"omitempty"`
	}

	func NewUser(dn string) *User {...}
	func (o *User) FilterObjectClass() string {...}
	func (o *User) Copy() crud.Item {...}
	func (o *User) Dn() string {...}
	func (o *User) MarshalLDAP() (*ldap.Entry, error) {...}
	func (o *User) UnmarshalLDAP(e *ldap.Entry) error {...}

For the conversion following rules are followed:

Names of are derived from the names defined in the schema files by first applying
strings.Title() and then removing every "-". Eg. "nick-name" would be converted to
"NickName".

The struct contains exported boolean fields for every object class it represents.

The struct contains every attribute of the object classes it represents as exported field.
If an attribute is marked as single-valued in the schema, the field has the type string.
Otherwise it has the type []string. If an attribute is used in the object class, but is not
defined in a read schema, it has also the type []string. No further distinctions are made so
far. Maybe at a later point, also the syntax definiton [4] of attributes will be considered,
but I don't know of any LDAP client that does this.

The generated methods are the following:

Methods to marshal and unmarshal from ldap.Entry. When marshalling,
MUST attributes are checked when the matching object class boolean field is true.

Copy method to create a copy of the struct.

Dn and SetDn methods to get and set the dn.

FilterObjectClass method to get the preferred object class for searching.

If two additional fields are present, a convinience method setting the DN to the current values is generated.
The field "DNFormat" field is a fmt.Sprintf format string. The field "DNAttributes" is a list of attribute names,
ordered for usage with the format string.

Example:
	[
	{
		"Name": "User",
		"Desc": "A user with several extensions",
		"ObjectClasses": ["posixAccount","person"],
		"FilterObjectClass": "posixAccount",
		"DNFormat":"sn=%v",
		"DNAttributes":["sn"]
	}
	]

Would generate this code:
	func (o *User) FormatDn() {
		o.dn = fmt.Sprintf("sn=%v,sn=%v", []string{o.Sn, o.Sn}...)
	}


[1] https://github.com/bytemine/ldap-crud/crud

[2] https://github.com/rbns/ldap

[3] https://raw.githubusercontent.com/rbns/ldap/master/entry.go

[4] https://www.ietf.org/rfc/rfc4517.txt

*/
package main
