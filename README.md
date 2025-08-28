# gmcwebapp
The GMC web app is a web application server designed to facilitate the
management of geologic material repositories.

## Configuration
Please see `gmc.sample.yaml` for configuration examples. The configuration
file is first searched for in the location specified at the command line
with `-conf`, then in `/etc/gmc.yaml`, and finally in `$HOME/.config/gmc.yaml`
on Linux. On Windows, the configuration file will be searched for in 
`C:\Users\<your username>\appdata\Roaming\gmc.yaml`

## HTTPS
HTTPS is optionally supported when the path to the listening certificate and key
are defined in the gmc.yaml file. If no key is certificate or key is defined, 
the server will default to serving connections over http. The certificate and 
key can be configured using the following fields in configuration yaml file:

* `listen_certificate`: The path to a user-provided SSL certificate.
* `listen_key`: The path to a user-provided SSL key.

## Database Configuration
Currently only PostgreSQL (v12+) is a supported database. Support for
additional databases can be added by implementing `gmc/db/DB`. The database 
connection can be configured with the following fields in the configuration 
yaml file: 

* `type`: The type of database, currently only 'postgres' is accepted.
* `url`: URL to database.
* `min_connections`: Minimum connections in connection pool, defaults to 0.
* `max_connections`: Maximum connections in connection pool, defaults to 4. 

## Search configuration
Currently only Elasticsearch (8.15+) is a supported search provider. Support 
for additional providers can be added by implementing `gmc/search/Search`. The 
connection to the search index requires the following fields in the 
configuration yaml file:

* `type`: The type of database, currently only 'elastic' is accepted.
* `url`: URL to search provider.
* `user`: Username for search provider access.
* `pass`: Password for search provider access.
* `index`: Name of search index. Defaults to "gmc".

## File Store Configuration
An array of redundant filestore locations can be defined. Both S3 and 
Directory-based file stores are supported. Support for additional file stores 
can be added by implementing `gmc/filestore/FileStore`. Every filestore 
defined will should have a `type` field. As of now, valid types include 
's3' and 'dir'.
 
* `type`: The type of filestore, currently 'dir' and 's3' are accepted.

### Directory filestore
Filestores can simply be a file directory accessible to the server, in which 
case a 'path' field will need to be defined in the configuration yaml file 
that defines the path to the filestore's folder. 

* `path`: Path to filestore folder.

### S3 filestore
Filestores can be stored on an S3-compatible object storage, in which case the 
following fields are required: 

* `endpoint`: URL to the S3 endpoint.
* `region`: Region of S3 service (this may be optional depending on S3 
  resource configuration).
* `secure`: Whether or not the connection will use HTTPS. Defaults to false.
* `accesskeyid`: ID for the endpoint's access key.
* `secretaccesskey`: Password for the endpoint's access key.
* `bucket`: Name of the bucket the filestore is stored in.

## Authentication Configuration
Four authentication schemata are supported: 'file', 'LDAP', 'token', and 
'always'. Support for additional authentication schemata can be added by 
implementing `gmc/auth/Auth`. The type of authentication is defined using the 
`type` field.

* `type`: The type of authentication schemata, currently 'file', 'LDAP', 
  'token', and 'always' are accepted.

### File Authentication
The file authentication schema uses a fixed file to authenticate users. The
file is cached on execution, so changes to the file require a restart. The
configuration yaml file requires a `path` field to define the file path.

* `path`: Path to authentication file.

### LDAP Authentication
The LDAP authentication schema establishes a connection to an LDAP server
to authenticate users. To configure the connection to the LDAP authentication
server, add the following values to the configuration yaml file based on the
LDAP server's configuration requirements:

* `name`: The authentication method name. Defaults to `ldap`.
* `ldap_url`: The LDAP server's URL.
* `bind_as_user`: A boolean flag that controls whether admin or user credentials
  are used to bind to the LDAP server. Defaults to false.
* `disable_certificate_verify`: A boolean flag to skip certificate verification
  when connecting to the LDAP server. Defaults to false.
* `bind_dn`: The distinguished name (DN) used to bind to the LDAP server when
 `bind_as_user` is false.
* `bind_password`: The password used to bind to the LDAP server when 
  `bind_as_user`
  is false.
* `base_dn`: The starting point for searches of the LDAP directory.
* `user_search`: The filter used to narrow searches of the LDAP directory.
* `ca`: The path to a user-provided certificate. If empty, the device's default
  certificates are used. `ca` is only read when disable_certificate_verify is
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

Example 3: This configuration uses a custom certificate with the LDAP server.  
The certificate field points to the certificate path.
```
- type: ldap
  ldap_url: ldaps://SOAADCDC06.soa.alaska.gov
  bind_as_user: true
  base_dn: ou=User,ou=akk,ou=tru,ou=akd,dc=poc,dc=wv,dc=com
  ca: "/home/user/my_cert.pem"
```

### Token Authentication
The token authentication schema uses the existing database connection to
authenticate tokens used by the companion GMC Android app. Installations
not using the companion Android app do not require token authentication.
Tokens can be managed via `gmc token`.

### Always Authentication
The always authentication accepts any username and password combination and
authenticates using whatever username is provided. It is primarily intended
for testing and is not appropriate for non-development installations. In the
configuration yaml file, setting the `allow` field to true or false (default 
to true) determines if all login attempts succeed or fail respectively.

* `allow`: Boolean field, determines if all login attempts succeed or fail. 

## Getting Started
Create a fresh PostgreSQL database owned by the user you intend to run the
application as. Configure the application using the instructions above and
`gmc.sample.yaml` as a reference. If you do not have existing data in your
database, execute `gmc db init` to initialize your database with the current
data model. Once a database is initialized and populated, run `gmc index` to 
index the database into the search provider. `gmc index` should be ran 
regularly to ensure that the data in the webapp is up-to-date. Start the 
server via `gmc start`.

## Development
Development requires a minimum of Go 1.25. For developers it
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
