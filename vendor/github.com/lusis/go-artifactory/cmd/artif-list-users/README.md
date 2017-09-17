# `artif-list-users`
```
usage: artif-list-users [<flags>]

List all users in Artifactory

Flags:
  --help              Show help (also see --help-long and --help-man).
  -F, --format=table  Format to show results [table, csv, list (usernames only - useful for piping)]
  --separator=","     separator for csv output
  --version           Show application version.
```

# Examples

## Default
`artif-list-users`

```
+-----------------------+------------------------------------------------------------------------------------------+
|         NAME          |                                           URI                                            |
+-----------------------+------------------------------------------------------------------------------------------+
| admin                 | http://artifactory/artifactory/api/security/users/admin                                  |
| deb                   | http://artifactory/artifactory/api/security/users/deb                                    |
| ci                    | http://artifactory/artifactory/api/security/users/ci                                     |
+-----------------------+------------------------------------------------------------------------------------------+
```

## csv
`artif-list-users -F csv`

```
admin,http://artifactory/artifactory/api/security/users/admin
dev,http://artifactory/artifactory/api/security/users/deb
ci,http://artifactory/artifactory/api/security/users/ci
```

Optionally you can provide a separator like so:

`artif-list-users -F csv --separator="|"` (pipe separated)

## list
Useful for piping into other scripts such as `artif-get-user`

`artif-list-users -F list`

```
admin
deb
ci
```
