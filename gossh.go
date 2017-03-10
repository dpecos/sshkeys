package main

import (
	"log"
)

func main() {

	hosts, acls, err := LoadConfig("config.ini")

	if err != nil {
		log.Panic(err)
	}

	err = UploadKeys("./id_rsa", "./keys", hosts, acls)

	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
}
