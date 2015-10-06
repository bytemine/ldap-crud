/*
Package slapd provides clean slapd instances for testing purposes.

Example usage:

	func TestLdap(t *testing.T) {
		// create a new slapd using the default config slapd.DefaultConfig
		s := slapd.New(nil)
		err := s.StartAndInitialize()
		if err != nil {
			t.Error(err)
		}
		defer s.Stop()

		// do some fancy stuff. the configuration used for this slapd instance can be accessed
		// at s.Config . E.g. the address is found at s.Config.Address() . If you need to access other
		// values, you can impement your own slapd.Configurer

	}

Note that the DefaultConfig only includes the core.schema of OpenLDAP, which is found in the "schema" subdirectory
of this package. The path to "schema" gets exported as Schemadir. If you want to use another set of schemata in
the DefaultConfig you can set it like this (note that the order matters for OpenLDAPs slapd!):

	slapd.DefaultConfig.Schemas = []string{"/my/schema/path/core.schema", ...}

*/
package slapd
