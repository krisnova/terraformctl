package terraformctl

import (
	"fmt"
	"github.com/kris-nova/kubicorn/cutil/logger"
	"google.golang.org/grpc"
)

// TerraformCtlClient is a wrapper for the gRPC client that we will use to interface with the terraformctl gRPC server.
// Use this instead of calling gRPC directly so that we can implement authentication and other features in the future.
type TerraformCtlClient struct {
	hostname string
	port     int
	Client   TerraformCTLAPIClient
}

// NewTerraformCtlClient will return a new terraformctl client that could be used to interface with a terraformctl gRPC server.
// hostname is the resolvable part of a URI that could be used to connect to a terraformctl gRPC server.
//      examples: localhost, www.mydomain.com, localdns.entry
// port is the TCP port to connect to the terraformctl gRPC server on
func NewTerraformCtlClient(hostname string, port int) *TerraformCtlClient {
	return &TerraformCtlClient{
		hostname: hostname,
		port:     port,
	}
}

// dialable returns a gRPC dialable string built from the hostname and port set are initilization of the TerraformCtlClient
func (c *TerraformCtlClient) dialable() string {
	return fmt.Sprintf("%s:%d", c.hostname, c.port)
}

// Connect will attempt to open a connection to a terraformctl gRPC server and will return a client based on the
// protobuf definition of the project.
func (c *TerraformCtlClient) Connect() error {
	logger.Info("Connecting to gRPC server [%s:%d]", c.hostname, c.port)
	connection, err := grpc.Dial(c.dialable(), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("Unable to connect to host [%s] with error message: %v", c.dialable(), err)
	}
	//defer connection.Close()
	client := NewTerraformCTLAPIClient(connection)
	c.Client = client
	return nil
}
