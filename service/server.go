package terraformctl

import "golang.org/x/net/context"

// Will ensure that TerraformCTLAPIServerImplementation is an implementation of TerraformCTLAPIServer
var _ TerraformCTLAPIServer = TerraformCTLAPIServerImplementation{}

// TerraformCTLAPIServerImplementation is an implementation of TerraformCTLAPIServer and is the struct
// that defines all interactions betweent the client and the server.
type TerraformCTLAPIServerImplementation struct {
}


// Get is used to retrieve data from a running gRPC service.
func (t TerraformCTLAPIServerImplementation) Get(context.Context, *GetRequest) (*GetResponse, error) {
	response := &GetResponse{}
	return response, nil
}

// Set is used to set new data to a running gRPC service. Set works like an upsert and will update or create
// as needed.
func (t TerraformCTLAPIServerImplementation) Set(context.Context, *SetRequest) (*SetResponse, error) {
	response := &SetResponse{}
	return response, nil
}
