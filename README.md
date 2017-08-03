Simple access control list for your servers infraestructure using just SSH keys.

About
=====

With sshkeys it's really easy to define your hosts list and who can access each of them. You can specify a list of users or use ACLs (named user lists) so a change in a team would be as easy as modifying that ACL and redeploying the SSH keys.

sshkeys works rewriting ~/.ssh/authorized_keys file of the remote hosts, adding only the required public keys ot the users allowed to access that host.

Important: in order to properly deploy to remote hosts, SSH access is required. The user executing sshkey must have access to remote hosts (with the account specified in the config file) in order to deploy the authorized keys. Make sure current user is always included in every deployment.

Installation
============

Golang is required to build sshkeys.

```
$ go get github.com/dpecos/sshkeys
$ go build github.com/dpecos/sshkeys
```

Usage
=====
* CONFIG is the path to the YML to use for SSH keys deployment
* PRIVATEKEY is the path to the private key to use for the hosts' account (the same for all of them)
* KEYRING is the path to the directory containing the users SSH public files, named as username.pub, where username is the name of the user specified in the config file

```
$ sshkeys --help
usage: sshkeys --config=CONFIG --privateKey=PRIVATEKEY --keyring=KEYRING [<flags>] <command> [<args> ...]

Deploy SSH keys to remote hosts

Flags:
  --help                   Show context-sensitive help (also try --help-long and --help-man).
  --config=CONFIG          Configuration file specifying hosts, grants and ACLs
  --privateKey=PRIVATEKEY  Private key to use to upload authorized_keys file to remote hosts
  --keyring=KEYRING        Path to a folder containing all user public keys (each one as a separate file named like
                           user.pub)

Commands:
  help [<command>...]
    Show help.

  deploy
    Generate and upload an authorized_keys file to remote hosts defined in configuation file
```

Example config file
===================

    [Hosts]
    server1.host = server1.test.com
    server1.account = root
    server1.acls = devops, developers
    server1.users = dpecos, john

    gw.host = gw.test.com
    gw.account = user
    gw.acls = devops

    private.host = priv.test.com
    private.account = dpecos
    private.users = dpecos

    [ACLs]
    devops = john
    developers = james, una
