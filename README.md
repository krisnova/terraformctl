# Terraformctl

Manage and mutate infrastructure with Terraform by running it in Kubernetes as a controller!

## Running

### Environmental Variables

`TERRAFORMCTL_HOSTNAME` can be used to override the hostname to use to connect to a listening gRPC server.
`TERRAFORMCTL_PORT` can be used to override the port to use to connect to a listening gRPC server.

## Developing



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
