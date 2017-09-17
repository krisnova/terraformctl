# go-artifactory
[![Build Status](https://travis-ci.org/lusis/go-artifactory.svg?branch=master)](https://travis-ci.org/lusis/go-artifactory)
Go library and utilities for interacting with [Artifactory](http://jfrog.com)

## Structure change notice
This was one of the first big golang libraries I wrote. I fought GOPATH for a long time. I finally "gave up" after having worked on go projects with other people.
This means that going forward, the library will now be much more approachable to other folks. 

To this end while the existing `src/artifactory.vXX` directories for previous versions are still available, the new model going forward is to shove the versions in the root directory of the repo (still versioned).
All utilities are now in the `cmd` directory and go-gettable.

This change should preserve existing usages externally (the only major one I know of is terraform's artifactory remote state provider).

I've also used this reorg to implement some linting standards which has cleaned up the code greatly.

Long story short, if you use this library you should upgrade to the new path for import (github.com/lusis/go-artifactory/artifactory.vXX) going forward.

## Breaking API changes due to move starting with `v51`
Additionally as of `v51`, the reorg and subsequent linting addition there were breaking API changes. These were mainly around stuttering and case changes on struct fields. The older versions are still valid and work in `v401` and `v491`.

## Usage
There are two ways to use this:
- as a library
- via the bundled utilities

### configuration
The following four environment variables are supported for configuring the artifactory client:

- `ARTIFACTORY_URL`
- `ARTIFACTORY_TOKEN`
- `ARTIFACTORY_PASSWORD`
- `ARTIFACTORY_USERNAME`

Newer version of artifactory support an API key that gets passed as a header instead of using basic auth. In the event that username, password and token are all specified, the token takes precedence.

Note that `ARTIFACTORY_URL` should be the base path to artifactory. This will be appended with the api paths e.g. `/api/security/users`). For this reason and due to issues with trailing slashes, if you have a trailing slash in the `ARTIFACTORY_URL` this will be trimmed.

You can read more about how artifactory authenticates the API [here](https://www.jfrog.com/confluence/display/RTF/Artifactory+REST+API)

### as a library
```go
package main

import (
	"fmt"

	artifactory "github.com/lusis/go-artifactory/artifactory.v51"
)

func main() {
	/*
		NewClientFromEnv requires two or three environment variables depending:
		- ARTIFACTORY_URL (i.e. https://myartifactory.domain.com/artifactory)

		and one of either:
		- ARTIFACTORY_TOKEN (this is the API key in newer versions of artifactory)

		or

		- ARTIFACTORY_USERNAME
    		- ARTIFACTORY_PASSWORD
	*/
	client := artifactory.NewClientFromEnv()
	data, err := client.GetUsers()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", data)
	}
}
```

### bundled utilities
All individual utilities are now go-gettable:

```
go get github.com/lusis/go-artifactory/cmd/artif-list-repos
```

or you can clone the entire repo and build them with the `Makefile`
```

git clone https://github.com/lusis/go-artifactory.git
cd go-artifactory
make all
ARTIFACTORY_URL=https://artifactory.domain.com/artifactory ARTIFACTORY_USERNAME=foo ARTIFACTORY_PASSWORD=bar bin/artif-list-repos
or
ARTIFACTORY_URL=https://artifactory.domain.com/artifactory ARTIFACTORY_TOKEN=XXXXXXX bin/artif-list-repos
```

```
+---------------------------------+---------+--------------------------------+------------------------------------------------------------------------+
|               KEY               |  TYPE   |          DESCRIPTION           |                                  URL                                   |
+---------------------------------+---------+--------------------------------+------------------------------------------------------------------------+
| dev-docker-local                | LOCAL   | Development Docker Registry    | https://artifactory/artifactory/dev-docker-local       |
| docker-local-v2                 | LOCAL   | Production Docker Repository   | https://artifactory/artifactory/docker-local-v2        |
| ext-release-local               | LOCAL   | Local repository for third     | https://artifactory/artifactory/ext-release-local      |
|                                 |         | party libraries                |                                                                        |
| ext-snapshot-local              | LOCAL   | Local repository for third     | https://artifactory/artifactory/ext-snapshot-local     |
|                                 |         | party snapshots                |                                                                        |
| libs-release-local              | LOCAL   | Local repository for in-house  | https://artifactory/artifactory/libs-release-local     |
|                                 |         | libraries                      |                                                                        |
| libs-snapshot-local             | LOCAL   | Local repository for in-house  | https://artifactory/artifactory/libs-snapshot-local    |
|                                 |         | snapshots                      |                                                                        |
| plugins-release-local           | LOCAL   | Local repository for plugins   | https://artifactory/artifactory/plugins-release-local  |
| plugins-snapshot-local          | LOCAL   | Local repository for plugins   | https://artifactory/artifactory/plugins-snapshot-local |
|                                 |         | snapshots                      |                                                                        |
| bower-remote                    | REMOTE  |                                | https://github.com/                                                    |
| docker-remote                   | REMOTE  |                                | https://registry-1.docker.io/                                          |
| jcenter                         | REMOTE  | Bintray Central Java           | http://jcenter.bintray.com                                             |
|                                 |         | repository                     |                                                                        |
| npm-remote                      | REMOTE  |                                | https://registry.npmjs.org                                             |
| nuget-remote                    | REMOTE  |                                | https://www.nuget.org/                                                 |
| pypi-remote                     | REMOTE  |                                | https://pypi.python.org                                                |
| rubygems-remote                 | REMOTE  |                                | https://rubygems.org/                                                  |
| nexus-releases-remote 	  | REMOTE  |                                | https://nexus/content/repositories/releases/        |
| libs-release                    | VIRTUAL |                                | https://artifactory/artifactory/libs-release           |
| libs-snapshot                   | VIRTUAL |                                | https://artifactory/artifactory/libs-snapshot          |
| plugins-release                 | VIRTUAL |                                | https://artifactory/artifactory/plugins-release        |
| plugins-snapshot                | VIRTUAL |                                | https://artifactory/artifactory/plugins-snapshot       |
| remote-repos                    | VIRTUAL |                                | https://artifactory/artifactory/remote-repos           |
+---------------------------------+---------+--------------------------------+------------------------------------------------------------------------+
```

## TODO
- Flesh out more tests
- More utilities
- Wrapper cli for sub-utilities

