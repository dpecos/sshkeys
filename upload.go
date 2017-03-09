package main

import (
	"fmt"
)

type AuthorizedKey struct {
	User    string
	Account string
	Host    string
}

func putKeysInHost(keysPath string, keys []AuthorizedKey) {
	for _, key := range keys {
		fmt.Printf("Uplading key for \"%s\" to \"%s@%s\"\n", key.User, key.Account, key.Host)
	}
}

func UploadKeys(keysPath string, hosts []Host, acls map[string]ACL) error {
	var keys []AuthorizedKey
	keysAdded := make(map[string]bool)

	for _, host := range hosts {
		for _, user := range host.Users {
			if _, seen := keysAdded[user+host.Account+host.Host]; !seen {
				keysAdded[user+host.Account+host.Host] = true
				keys = append(keys, AuthorizedKey{user, host.Account, host.Host})
			}
		}
		for _, acl := range host.ACLs {
			for _, user := range acl.Users {
				if _, seen := keysAdded[user+host.Account+host.Host]; !seen {
					keysAdded[user+host.Account+host.Host] = true
					keys = append(keys, AuthorizedKey{user, host.Account, host.Host})
				}
			}
		}
	}

	putKeysInHost(keysPath, keys)

	return nil
}
