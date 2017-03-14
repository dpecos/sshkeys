package main

import (
	"log"
	"os"

	"github.com/dpecos/sshkeys/config"
	"github.com/dpecos/sshkeys/utils"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app        = kingpin.New("sshkeys", "Deploy SSH keys to remote hosts")
	configFile = app.Flag("config", "Configuration file specifying hosts, grants and ACLs").Required().String()
	privateKey = app.Flag("privateKey", "Private key to use to upload authorized_keys file to remote hosts").Required().String()
	keyring    = app.Flag("keyring", "Path to a folder containing all user public keys (each one as a separate file named like user.pub)").Required().String()
	deploy     = app.Command("deploy", "Generate and upload an authorized_keys file to remote hosts defined in configuation file")
)

func main() {

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {

	case deploy.FullCommand():
		hosts, acls, err := config.LoadConfig(*configFile)
		if err != nil {
			log.Panic(err)
		}

		if err := utils.UploadKeys(*privateKey, *keyring, hosts, acls); err != nil {
			log.Fatalf("ERROR: %s\n", err)
		}
	}

}
