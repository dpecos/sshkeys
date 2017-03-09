package main

func main() {

	hosts, acls, err := LoadConfig("config.ini")

	if err != nil {
		panic(err)
	}

	err = CheckKeys("./keys", hosts, acls)

	if err != nil {
		panic(err)
	}

	err = UploadKeys("./keys", hosts, acls)

	if err != nil {
		panic(err)
	}
}
