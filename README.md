## Configuration
Please see `gmc.sample.yaml` for configuration examples. The configuration
file is first searched for in the location specified at the command line
with `-conf`, then in `/etc/gmc.yaml`, and finally in `$HOME/.config/gmc.yaml`
(on Linux.) Configuration requires at minimum a database source and a
filestore. Neither is checked on startup, so fillers may be used for
application testing and development.

#### Database
Currently only PostgreSQL (v12+) is a supported database. Support for
additional databases can be added by implementing `gmc/db/DB`.

#### File Store
Both S3 and Directory-based file stores are supported. Support for additional
file stores can be added by implementing `gmc/filestore/FileStore`.

#### Authentication
Three authentication schemata are supported: file, LDAP, and token.  The file
authentication schema uses a fixed file to authenticate users. The file is
cached on execution, so changes to the file require a restart.  The LDAP
authentication schema connects to a LDAP serve and can be configured in
`/etc/gmc.yaml`.  The token authentication schema uses a database table to
authenticate users.  The database URL is specified in `/etc/gmc.yaml`. Support
for additional file stores can be added by implementing `gmc/auth/Auth`.

## Development
Development requires a minimum of Go 1.22. For developers it
is recommended to configure `auto_shutdown` to true, and to start the
application (on Linux) from the project directory with
`while true; do ./gmc -assets assets start; done` - this will enable
dynamic reloading of HTML, CSS, templates, and queries, while also
automatically restarting the application on recompilation.
