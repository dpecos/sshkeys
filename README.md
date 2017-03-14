Example config file
===================

    [Hosts]
    server1.host = server1.test.com
    server1.account = root
    server1.acls = devops, developers
    server1.users = dpecos

    gw.host = gw.test.com
    gw.account = user
    gw.acls = devops

    private.host = priv.test.com
    private.account = dpecos
    private.users = dpecos

    [ACLs]
    devops = john
    developers = james, una