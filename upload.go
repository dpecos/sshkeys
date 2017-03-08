package main

import (
	"fmt"
)

func putKeyInHost(user string, account string, host string) {
	fmt.Printf("Uplading key for \"%s\" to \"%s@%s\"\n", user, account, host)
}

func UploadKeys(hosts []Host, acls map[string]ACL) error {
	for _, host := range hosts {
		for _, user := range host.Users {
			putKeyInHost(user, host.Account, host.Host)
		}
		for _, acl := range host.ACLs {
			for _, user := range acl.Users {
				putKeyInHost(user, host.Account, host.Host)
			}
		}
	}
	return nil
}
