# Terraformctl

Manage and mutate infrastructure with Terraform by running it in Kubernetes as a controller!

## About

This is an example of infrastructure as cloud native software.
This repository is not intended to be used in production, but rather offers a starting point for people to start looking at what it would take to run infrastructure as software.

More information can be found on my blog [here](http://www.nivenly.com/i-ran-terraform-in-kubernetes/)

## Running

### Environmental Variables

`TERRAFORMCTL_HOSTNAME` can be used to override the hostname to use to connect to a listening gRPC server.
`TERRAFORMCTL_PORT` can be used to override the port to use to connect to a listening gRPC server.

## Developing

## Building and pushing

Sorry but I hard coded everything for a demo.. be ready to hack the Makefile (please open a PR if you want!)

```
make build push deploy
```

This also assumes you have Kubernetes up and running already.

### Working with the gRPC definitions

You will need to have `protoc` and `grpc` installed.

```bash
go get google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
```

Now you can make changes to `service/terraformctl.proto` and run the following command to update the plugin.

```bash
make proto
```

A change to the gRPC might be needed in `service/server.go` if the new gRPC expects new logic.
