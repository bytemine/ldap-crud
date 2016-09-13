/*
Package slapd provides clean slapd instances for testing purposes.

You can use it on two ways:

Create a Config, and let the Configure() method create a exec.Command for you to run.
This can at times be tricky, because there is no clean way to test if the slapd is already
accepting connections.

Manually configure and run slapd and create a Config matching the values where appropriate.
You can then use the Clean() and Initialize() methods to reset the contents in the ldap. 
This is much more reliable.
*/
package slapd
