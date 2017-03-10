package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
)

type AuthorizedKey struct {
	User    string
	Account string
	Host    string
}

func readKeyFile(keyspath string, user string) (string, error) {
	keyname := filepath.Join(keyspath, user+".pub")
	contents, err := ioutil.ReadFile(keyname)
	return string(contents), err
}

func createAuthorizedKeysFile(keysPath string, keys []AuthorizedKey) (string, error) {
	file, err := ioutil.TempFile("", "auth")
	if err != nil {
		return "", err
	}
	defer file.Close()

	for _, key := range keys {
		log.Printf("Adding key for \"%s\" to \"%s@%s\"\n", key.User, key.Account, key.Host)
		keyContents, err := readKeyFile(keysPath, key.User)
		if err != nil {
			defer os.Remove(file.Name())
			return "", err
		}
		_, err = file.WriteString(keyContents + "\n")
		if err != nil {
			defer os.Remove(file.Name())
			return "", err
		}
	}

	return file.Name(), nil
}

func uploadFileToHost(privateKey string, host Host, file string, remoteFilename string) error {
	clientConfig, _ := auth.PrivateKey(host.Account, privateKey)
	client := scp.NewClient(host.Host+":22", &clientConfig)

	err := client.Connect()
	if err != nil {
		log.Fatal("Couldn't establisch a connection to the remote server ", err)
		return err
	}
	defer client.Session.Close()

	reader, _ := os.Open(file)
	defer reader.Close()

	client.CopyFile(reader, remoteFilename, "0644")

	return nil
}

func putKeysInHost(privateKey string, keysPath string, keys []AuthorizedKey, host Host) error {
	file, err := createAuthorizedKeysFile(keysPath, keys)
	if err != nil {
		return err
	}
	defer os.Remove(file)

	err = uploadFileToHost(privateKey, host, file, "~/.ssh/authorized_keys")

	return err
}

func UploadKeys(privateKey, keysPath string, hosts []Host, acls map[string]ACL) error {
	for _, host := range hosts {
		var keys []AuthorizedKey
		keysAdded := make(map[string]bool)

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

		if len(keys) > 0 {
			log.Printf("Uploading keys to %s@%s...\n", host.Account, host.Host)
			err := putKeysInHost(privateKey, keysPath, keys, host)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
