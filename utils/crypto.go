package utils

import (
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
		if _, err = file.WriteString(keyContents + "\n"); err != nil {
			defer os.Remove(file.Name())
			return "", err
		}
	}

	return file.Name(), nil
}
