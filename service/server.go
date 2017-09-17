package terraformctl

import (
	"fmt"
	"github.com/kris-nova/terraformctl/parser"
	"github.com/kris-nova/terraformctl/storage"
	"golang.org/x/net/context"
)

// Will ensure that TerraformCTLAPIServerImplementation is an implementation of TerraformCTLAPIServer
var _ TerraformCTLAPIServer = TerraformCTLAPIServerImplementation{}

// TerraformCTLAPIServerImplementation is an implementation of TerraformCTLAPIServer and is the struct
// that defines all interactions betweent the client and the server.
type TerraformCTLAPIServerImplementation struct {
	cacher    storage.TerraformCtlCacher
	persister storage.TerraformCtlPersister
}

// NewTerraformCTLAPIServerImplementation will initialize a new TerraformCTLAPIServerImplementation that will be called via the gRPC protocol.
func NewTerraformCTLAPIServerImplementation(cacher storage.TerraformCtlCacher, persister storage.TerraformCtlPersister) *TerraformCTLAPIServerImplementation {
	return &TerraformCTLAPIServerImplementation{
		cacher:    cacher,
		persister: persister,
	}
}

// Get is used to retrieve data from a running gRPC service.
func (t TerraformCTLAPIServerImplementation) Get(ctx context.Context, request *GetRequest) (*GetResponse, error) {
	response := &GetResponse{}
	t.cacher.Lock()
	config, err := t.cacher.Get(request.Name)
	t.cacher.Unlock()
	if err != nil {
		return nil, fmt.Errorf("Unable to get configuration [%s] with error message: %v", request.Name, err)
	}
	str, err := config.String()
	if err != nil {
		return nil, fmt.Errorf("Unable to parse config as string with error: %v", err)
	}
	response.Config = str
	response.Name = request.Name
	return response, nil
}

// Set is used to set new data to a running gRPC service. Set works like an upsert and will update or create
// as needed.
func (t TerraformCTLAPIServerImplementation) Set(ctx context.Context, request *SetRequest) (*SetResponse, error) {
	response := &SetResponse{}
	configStr, err := parser.NewTerraformConfigurationFromEncodedString(request.Config)
	if err != nil {
		return nil, fmt.Errorf("Unable to generate new Terraform configuration with error: %v", err)
	}
	t.cacher.Lock()
	err = t.cacher.Set(request.Name, configStr)
	t.cacher.Unlock()
	if err != nil {
		return nil, fmt.Errorf("Unable to set configuration [%s] with error: %v", request.Name, err)
	}
	response.Config = request.Config
	response.Name = request.Name
	return response, nil
}
