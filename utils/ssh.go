package utils

import (
	"log"
	"os"

	config "github.com/dpecos/sshkeys/config"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
)

func uploadFileToHost(privateKey string, host config.Host, file string, remoteFilename string) error {
	clientConfig, _ := auth.PrivateKey(host.Account, privateKey)
	client := scp.NewClient(host.Host+":22", &clientConfig)

	err := client.Connect()
	if err != nil {
		log.Fatal("Couldn't establisc a connection to the remote server: ", err)
		return err
	}
	defer client.Session.Close()

	reader, _ := os.Open(file)
	defer reader.Close()

	client.CopyFile(reader, remoteFilename, "0644")

	return nil
}

func putKeysInHost(privateKey string, keysPath string, keys []AuthorizedKey, host config.Host) error {
	file, err := createAuthorizedKeysFile(keysPath, keys)
	if err != nil {
		return err
	}
	defer os.Remove(file)

	err = uploadFileToHost(privateKey, host, file, "~/.ssh/authorized_keys")

	return err
}

func UploadKeys(privateKey, keysPath string, hosts []config.Host, acls map[string]config.ACL) error {
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
