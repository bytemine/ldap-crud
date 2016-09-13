package slapd

import (
	"log"
)

func Example() {
	cmd, err := DefaultConfig.Configure()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	// Add basedn and rootdn entries.
	err = DefaultConfig.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	// Do something in the ldap...

	// Remove all you've done from the ldap.
	err = DefaultConfig.Clean()
	if err != nil {
		log.Fatal(err)
	}

	// Re-add basedn and rootdn.
	err = DefaultConfig.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	// Somehow reliably kill the slapd.

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	// Remove stuff created by Configure
	err = DefaultConfig.Unconfigure()
	if err != nil {
		log.Fatal(err)
	}
}

func Example_manual() {
	// Manually start a slapd configured like DefaultConfig (or your own).

	// Add basedn and rootdn entries.
	err = DefaultConfig.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	// Do something in the ldap...

	// Remove all you've done from the ldap.
	err = DefaultConfig.Clean()
	if err != nil {
		log.Fatal(err)
	}

	// Re-add basedn and rootdn.
	err = DefaultConfig.Initialize()
	if err != nil {
		log.Fatal(err)
	}
}
