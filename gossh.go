package main

import (
	"fmt"
)

func main() {
	acls, hosts, err := LoadConfig("config.ini")

	if err != nil {
		panic(err)
	}

	fmt.Printf("ACLs: %v", acls)
	fmt.Printf("Hosts: %v", hosts)
}
