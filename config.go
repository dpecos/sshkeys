package main

import (
	"fmt"
	"strings"

	ini "gopkg.in/ini.v1"
)

type Host struct {
	ID    string
	Host  string
	Admin string
	ACLs  []ACL
	Users []string
}

type ACL struct {
	Name  string
	Users []string
}

func loadACLs(cfg *ini.File) (map[string]ACL, error) {
	acls := make(map[string]ACL)
	section, err := cfg.GetSection("ACLs")
	if err != nil {
		return nil, fmt.Errorf("Section ACLs not found in config file: %s", err)
	}
	for _, key := range section.Keys() {
		acl := new(ACL)
		acl.Name = key.Name()
		acl.Users = key.Strings(",")
		acls[acl.Name] = *acl
	}
	return acls, nil
}

func getKey(section *ini.Section, key string) *ini.Key {
	val, err := section.GetKey(key)
	if err != nil {
		panic(fmt.Sprintf("Key %s not found", key))
	}
	return val
}

func loadHosts(cfg *ini.File, acls map[string]ACL) ([]Host, error) {
	var hosts []Host
	section, err := cfg.GetSection("Hosts")
	if err != nil {
		return nil, fmt.Errorf("Section Hosts not found in config file: %s", err)
	}
	for _, key := range section.Keys() {
		if strings.HasSuffix(key.Name(), ".host") {
			id := strings.Split(key.Name(), ".")[0]
			host := new(Host)
			host.ID = id
			host.Admin = getKey(section, id+".admin").String()
			host.Users = getKey(section, id+".users").Strings(",")
			aclNames := getKey(section, id+".acls").Strings(",")
			for _, name := range aclNames {
				host.ACLs = append(host.ACLs, acls[name])
			}
			hosts = append(hosts, *host)
		}
	}
	return hosts, nil
}

func LoadConfig(f string) (map[string]ACL, []Host, error) {
	cfg, err := ini.Load(f)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not open config file %s: %s", f, err)
	}

	acls, err := loadACLs(cfg)
	if err != nil {
		return nil, nil, err
	}

	hosts, err := loadHosts(cfg, acls)

	return acls, hosts, nil
}
