# gmcwebapp
The GMC web app is a web application server designed to facilitate the
management of geologic material repositories.

## Configuration
Please see `gmc.sample.yaml` for configuration examples. The configuration
file is first searched for in the location specified at the command line
with `-conf`, then in `/etc/gmc.yaml`, and finally in `$HOME/.config/gmc.yaml`
(on Linux.)

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
The LDAP authentication schema connects to an LDAP server.   To configure the LDAP authentication server add the following values to the configuration yaml file:
  - type: ldap
    ldap_url: ldap_url
    bind_dn: bind_dn
    bind_password: mypassword
    bind_as_user: true
    disable_certificate_verify: false
    base_dn: base_dn
    user_search: search_filter

Two common configurations using placeholder values are

	Example 1: bind_as_user = true
	  - type: ldap
			ldap_url: ldaps://ldap4.example.com
			bind_dn: uid=<BIND_USERNAME>,ou=User,o=co.wv.uk.eu

	Example 2: bind_as_user = false
		- type: ldap
			ldap_url: ldaps://ldap5.example.com
			bind_dn: uid=<BIND_USERNAME>,ou=User,o=co.wv.uk.eu
			bind_password: <BIND_PASSWORD>
			bind_as_user: false
			base_dn: ou=User,ou=akk,ou=tru,ou=akd,dc=poc,dc=wv,dc=com
			user_search: (CN={{.}})

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
