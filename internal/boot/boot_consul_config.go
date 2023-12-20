package boot

import (
	"github.com/fulltimelink/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/contrib/config/consul/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/hashicorp/consul/api"
)

type bootConfig struct {
	Address string
	Schema  string
	Path    string
}

func (b *bootConfig) Run() (config.Config, conf.Bootstrap) {
	consulClient, err := api.NewClient(&api.Config{
		Address: b.Address,
		Scheme:  b.Schema,
	})
	if err != nil {
		panic(err)
	}
	cs, err := consul.New(consulClient, consul.WithPath(b.Path))
	// consul中需要标注文件后缀，kratos读取配置需要适配文件后缀
	// The file suffix needs to be marked, and kratos needs to adapt the file suffix to read the configuration.
	if err != nil {
		panic(err)
	}
	c := config.New(config.WithSource(cs))

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	return c, bc
}

func NewBootConfig(address string, schema string, path string) *bootConfig {
	return &bootConfig{
		Address: address,
		Schema:  schema,
		Path:    path,
	}
}
