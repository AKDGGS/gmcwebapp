## Configuration
Please see `gmc.yaml` for configuration examples. `gmc.yaml` is searched for
in `/etc`, as well as `$HOME/.config` (on Linux), and finally in a location
specified at the command line with `-conf`.

Configuration requires a database source and a filestore. Neither is checked
on startup, so fillers may be used for application testing and development.

#### Database
Currently only PostgreSQL is a supported database. Support for additional
databases can be added by implementing `gmc/DB/DB`.

#### File Store
Both S3 and Directory-based file stores are supported. It is highly recommended
that S3 be used in production environments, as it offers a higher degree
of security and built-in support for ETags. Directory-based file stores are
intended for development and testing.


## Development
Development requires a minimum of Go 1.16 as go:embed is used extensively.
Building is as simple as `go build` in the project directory. For developers
it is recommended to start the application on Linux from the project directory
with `while true; do ./gmc start -assets assets -s; done` - this will enable
dynamic reloading of HTML, CSS, templates, and queries, while also
automatically restarting the application on recompilation.
