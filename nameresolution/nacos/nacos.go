package nacos

import (
	nr "github.com/dapr/components-contrib/nameresolution"
	"github.com/dapr/kit/logger"
	nacos_client "github.com/nacos-group/nacos-sdk-go/clients/nacos_client"
)

const daprMeta string = "DAPR_PORT" // default key for DAPR_PORT metadata

type client struct {
	*nacos_client.NacosClient
}

// NewResolver creates Consul name resolver.
func NewResolver(logger logger.Logger) nr.Resolver {
	return newResolver(logger, resolverConfig{}, &client{})
}

func newResolver(logger logger.Logger, resolverConfig resolverConfig, client clientInterface) nr.Resolver {
	return &resolver{
		logger: logger,
		config: resolverConfig,
		client: client,
	}
}
