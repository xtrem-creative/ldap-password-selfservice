# LDAP self-service password change

A minimalist LDAP self-service password change app, written in Go, using Buffalo.

## Configuration
You are required to set the following ENV variables to make it work:

```bash
# The LDAP address (with its port)
SELFSERVICE_LDAP_URL="ldap.mydomain.tld:port"
# The LDAP server protocol method
SELFSERVICE_LDAP_METHOD="tls" # plain or tls
# The LDAP bind DN
SELFSERVICE_LDAP_BIND_DN="cn=webuser,dc=mydomain,dc=tld"
# The LDAP bind password
SELFSERVICE_LDAP_PASSWORD="my_bind_dn_password"
# The LDAP root, where to find users
SELFSERVICE_LDAP_BASE="ou=People,dc=mydomain,dc=tld"
# A LDAP valid filter (can be set to empty)
SELFSERVICE_LDAP_FILTER="(&(objectClass=inetOrgPerson)(memberOf=cn=active,ou=Groups,dc=mydomain,dc=tld))"
```

Either set this config using a .env file in the same directory as the app, or set it using system utilities.

[Powered by Buffalo](http://gobuffalo.io)
