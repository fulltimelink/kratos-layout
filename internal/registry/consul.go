package registry

import (
	"github.com/fulltimelink/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/hashicorp/consul/api"
)

func NewConsulRegistry(c *conf.Registry) *consul.Registry {
	client, err := api.NewClient(&api.Config{
		Address: c.Consul.Address,
		Scheme:  c.Consul.Schema})
	if err != nil {
		panic(err)
	}
	return consul.New(client, consul.WithDeregisterCriticalServiceAfter(int(c.Consul.DeregisterCriticalServiceAfter)))
}
