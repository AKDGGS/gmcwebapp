# gmcwebapp
The GMC web app is a web application server designed to facilitate the
management of geologic material repositories.

## Configuration
Please see `gmc.sample.yaml` for configuration examples. The configuration
file is first searched for in the location specified at the command line
with `-conf`, then in `/etc/gmc.yaml`, and finally in `$HOME/.config/gmc.yaml`
(on Linux.)

## HTTPS
HTTPS is optionally supported when the path to the listening certificate and key
are defined in the gmc.yaml file.

### Database Configuration
Currently only PostgreSQL (v12+) is a supported database. Support for
additional databases can be added by implementing `gmc/db/DB`.

### File Store Configuration
Both S3 and Directory-based file stores are supported. Support for additional
file stores can be added by implementing `gmc/filestore/FileStore`.

### Authentication Configuration
Four authentication schemata are supported: file, LDAP, token, and always.
Support for additional authentication schemata can be added by implementing
`gmc/auth/Auth`.

#### File Authentication
The file authentication schema uses a fixed file to authenticate users. The
file is cached on execution, so changes to the file require a restart.

#### LDAP Authentication
The LDAP authentication schema establishes a connection to an LDAP server
to authenticate users.  To configure the connection to the LDAP authentication
server, add the following values to the configuration yaml file based on the
LDAP server's configuration requirements:

* `name`: The authentication method name. Defaults to 'ldap'.
* `ldap_url`: The LDAP server's URL.
* `bind_as_user`: A boolean flag that controls whether admin or user credentials
  are used to bind to the LDAP server.  Defaults to false.
* `disable_certificate_verify`: A boolean flag to skip certificate verification
  when connecting to the LDAP server.  Defaults to false.
* `bind_dn`: The distinguished name (DN) used to bind to the LDAP server when
  `bind_as_user` is false.
* `bind_password`: The password used to bind to the LDAP server when `bind_as_user`
  is false.
* `base_dn`: The starting point for searches of the LDAP directory.
* `user_search`: The filter used to narrow searches of the LDAP directory.
* `ca`: The path to a user-provided certificate,  If empty, the device's default
  certificates are used.  `ca` is only read when disable_certificate_verify is
  false.

Three common configurations using placeholder values are:

Example 1: This configuration uses the user's credentials to bind to the LDAP
server if bind_as_user is set to true.
```
- type: ldap
  ldap_url: ldaps://ldap4.example.com
  bind_as_user: true
  bind_dn: uid=<BIND_USERNAME>,ou=User,o=co.wv.uk.eu
```

Example 2: This configuration uses a specific user's credentials to bind to the
LDAP server.
```
- type: ldap
  ldap_url: ldaps://ldap4.example.com
  bind_dn: uid=<BIND_USERNAME>,ou=User,o=co.wv.uk.eu
  bind_password: <BIND_PASSWORD>
  base_dn: ou=User,ou=akk,ou=tru,ou=akd,dc=poc,dc=wv,dc=com
  user_search: (CN={{.}})
```

Example 3: This configuration uses a custom certificate with the LDAP server. The
certificate field points to the certificate path.
```
- type: ldap
  ldap_url: ldaps://SOAADCDC06.soa.alaska.gov
  bind_as_user: true
  base_dn: ou=User,ou=akk,ou=tru,ou=akd,dc=poc,dc=wv,dc=com
  ca: "/home/user/my_cert.pem"
```

#### Token Authentication
The token authentication schema uses the existing database connection to
authenticate tokens used by the companion GMC Android app. Installations
not using the companion Android app do not require token authentication.
Tokens can be managed via `gmc token`.

#### Always Authentication
The always authentication accepts any username and password combination and
authenticates using whatever username is provided. It is primarily intended
for testing and is not appropriate for non-development installations.

## Getting Started
Create a fresh PostgreSQL database owned by the user you intend to run the
application as. Configure the application using the instructions above and
`gmc.sample.yaml` as a reference. If you do not have existing data in your
database, execute `gmc db init` to initialize your database with the current
data model. Start the server via `gmc start`.

## Development
Development requires a minimum of Go 1.22. For developers it
is recommended to configure `auto_shutdown` to true, and to start the
application (on Linux) from the project directory with
`while true; do ./gmc -assets assets start; done` - this will enable
dynamic reloading of HTML, CSS, templates, and queries, while also
automatically restarting the application on recompilation.

### Customizing
The GMC app can be invoked with the assets parameter pointing to a directory,
this directory is then used instead of the built-in assets. If a specific
asset does not exist in the directory, it will fall back to the built-in
asset of the same name. As a result, installations can be customized
visually with templates (`assets/css`, `assets/tmpl`) and extended with
specialized Q/A reports (`assets/pg/qa`).
