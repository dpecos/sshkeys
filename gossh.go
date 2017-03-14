package main

import (
	"log"
	"os"

	"github.com/dpecos/gossh/config"
	"github.com/dpecos/gossh/utils"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app        = kingpin.New("gossh", "")
	configFile = app.Flag("config", "Configuration file specifying hosts, grants and ACLs").Required().String()
	privateKey = app.Flag("privateKey", "Private key to use to upload authorized_keys file to remote hosts").Required().String()
	keyring    = app.Flag("keyring", "Path to a folder containing all user public keys (each one as a separate file named like user.pub)").Required().String()
)

func main() {

	kingpin.MustParse(app.Parse(os.Args[1:]))

	hosts, acls, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Panic(err)
	}

	if err := utils.UploadKeys(*privateKey, *keyring, hosts, acls); err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
}
