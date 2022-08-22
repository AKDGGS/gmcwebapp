## Configuration
Please see `gmc.yaml` for configuration examples. `gmc.yaml` is searched for in
`/etc`, as well as `$HOME/.config` (on Linux), and finally in a location
specified at the command line with `-conf`.

Configuration requires at minimum a database source and a filestore. Neither is
checked on startup, so fillers may be used for application testing and
development.

#### Database
Currently only PostgreSQL (v12+) is a supported database. Support for
additional databases can be added by implementing `gmc/db/DB`.

#### File Store
Both S3 and Directory-based file stores are supported. Support for additional
file stores can be added by implementing `gmc/filestore/FileStore`.

#### Authentication
Two authentication schemes are supported, file and always. The file
authentication scheme uses a fixed file to authenticate users. The file is
cached on execution, so changes to the file require a restart. The always
scheme is intended for development and testing. Any attempt to authenticate
will immediately return success with the configured username. Support for
additional file stores can be added by implementing `gmc/auth/Auth`.

## Development
Development requires a minimum of Go 1.16 as go:embed is used extensively.
Building is as simple as `go build` in the project directory. For developers it
is recommended to configure `auto_shutdown` to true, and to start the
application (on Linux) from the project directory with
`while true; do ./gmc -assets assets start; done` - this will enable
dynamic reloading of HTML, CSS, templates, and queries, while also
automatically restarting the application on recompilation.
