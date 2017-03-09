package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type AuthorizedKey struct {
	User    string
	Account string
	Host    string
}

func readKeyFile(keyspath string, user string) (string, error) {
	keyname := filepath.Join(keyspath, user+".pub")
	file, err := os.OpenFile(keyname, os.O_RDONLY, 0444)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var text = make([]byte, 1024)
	for {
		n, err := file.Read(text)
		if err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}
	}

	return string(text), nil
}

func createAuthorizedKeysFile(keysPath string, keys []AuthorizedKey) (*os.File, error) {
	file, err := ioutil.TempFile("", "auth")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	for _, key := range keys {
		log.Printf("Adding key for \"%s\" to \"%s@%s\"\n", key.User, key.Account, key.Host)
		keyContents, err := readKeyFile(keysPath, key.User)
		if err != nil {
			defer os.Remove(file.Name())
			return nil, err
		}
		_, err = file.WriteString(keyContents + "\n")
		if err != nil {
			defer os.Remove(file.Name())
			return nil, err
		}
	}

	return file, nil
}

func uploadFileToHost(file *os.File, host Host) error {
	return nil
}

func putKeysInHost(keysPath string, keys []AuthorizedKey, host Host) error {
	file, err := createAuthorizedKeysFile(keysPath, keys)
	if err != nil {
		return err
	}
	print(file.Name())
	defer os.Remove(file.Name())

	err = uploadFileToHost(file, host)

	return err
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

		err := putKeysInHost(keysPath, keys, host)
		if err != nil {
			return err
		}
	}

	return nil
}
