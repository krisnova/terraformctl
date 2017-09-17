gomanta
=======

The gomanta package enables Go programs to interact with the Joyent Manta service.

[![wercker status](https://app.wercker.com/status/e315959b9089bba1da9aff10d36e5b4c/s/master "wercker status")](https://app.wercker.com/project/byKey/e315959b9089bba1da9aff10d36e5b4c)

## Installation

Use `go-get` to install gomanta
```
go get github.com/joyent/gomanta
```

## Packages

The gomanta package is structured as follow:

	- gomanta/localservices. This package provides local services to be used for testing.
	- gomanta/manta. This package interacts with the Manta API (http://apidocs.joyent.com/manta/).


## Documentation

Documentation can be found on godoc.

- [github.com/joyent/gomanta](http://godoc.org/github.com/joyent/gomanta)
- [github.com/joyent/gomanta/localservices](http://godoc.org/github.com/joyent/gomanta/localservices)
- [github.com/joyent/gomanta/manta](http://godoc.org/github.com/joyent/gomanta/manta)

## Contributing

Report bugs and request features using [GitHub Issues](https://github.com/joyent/gomanta/issues), or contribute code via a [GitHub Pull Request](https://github.com/joyent/gomanta/pulls). Changes will be code reviewed before merging. In the near future, automated tests will be run, but in the meantime please `go fmt`, `go lint`, and test all contributions.


## Developing

This library assumes a Go development environment setup based on [How to Write Go Code](https://golang.org/doc/code.html). Your GOPATH environment variable should be pointed at your workspace directory.

You can now use `go get github.com/joyent/gomanta` to install the repository to the correct location, but if you are intending on contributing back a change you may want to consider cloning the repository via git yourself. This way you can have a single source tree for all Joyent Go projects with each repo having two remotes -- your own fork on GitHub and the upstream origin.

For example if your GOPATH is `~/src/joyent/go` and you're working on multiple repos then that directory tree might look like:

```
~/src/joyent/go/
|_ pkg/
|_ src/
   |_ github.com
      |_ joyent
         |_ gocommon
         |_ gomanta
         |_ gosdc
         |_ gosign
```

### Recommended Setup

```
$ mkdir -p ${GOPATH}/src/github.com/joyent
$ cd ${GOPATH}/src/github.com/joyent
$ git clone git@github.com:<yourname>/gomanta.git

# fetch dependencies
$ git clone git@github.com:<yourname>/gocommon.git
$ git clone git@github.com:<yourname>/gosign.git
$ go get -v -t ./...

# add upstream remote
$ cd gomanta
$ git remote add upstream git@github.com:joyent/gomanta.git
$ git remote -v
origin  git@github.com:<yourname>/gomanta.git (fetch)
origin  git@github.com:<yourname>/gomanta.git (push)
upstream        git@github.com:joyent/gomanta.git (fetch)
upstream        git@github.com:joyent/gomanta.git (push)
```

### Run Tests

```
cd ${GOPATH}/src/github.com/joyent/gomanta
go test ./...
```

The `manta` package tests can also be run against live Manta. If you want to run this package, you can pass the `-live` flag and the `-key.name` flag (the latter is optional and defaults to `~/.ssh/id_rsa`) as shown below. Note that you can only run the `manta` package tests this way and running the rest of the test suite with these flags will result in a test runner error ("flag provided but not defined").

```
cd ${GOPATH}/src/github.com/joyent/gomanta
go test ./manta -live -key.name=~/.ssh/my_key
```


### Build the Library

```
cd ${GOPATH}/src/github.com/joyent/gomanta
go build ./...
```

## License
Licensed under [MPLv2](LICENSE).

Copyright (c) 2016 Joyent Inc.
Written by Daniele Stroppa <daniele.stroppa@joyent.com>
