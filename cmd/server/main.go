package main

import (
	"flag"
	"github.com/fulltimelink/kratos-layout/internal/boot"
	"github.com/fulltimelink/kratos-layout/internal/conf"
	"github.com/fulltimelink/kratos-layout/internal/server"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"os"

	"github.com/go-kratos/kratos/v2/transport"

	"github.com/go-kratos/kratos/v2/config"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

var (
	// --  @# address is config's address
	address string

	// --  @# schema is config's schema
	schema string

	// --  @# path is config's path
	path string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&address, "address", "127.0.0.1:8500", "config address, eg: -address 127.0.0.1:8500")
	flag.StringVar(&schema, "schema", "http", "config schema, eg: -schema http")
	flag.StringVar(&path, "path", "config/huidu-server/config.yaml", "config path, eg: -path conf/huidu-server/config.yaml")
}

func newApp(conf *conf.Server, r *consul.Registry, logger log.Logger, gs *grpc.Server, hs *http.Server, ps *server.PProfServer) *kratos.App {
	// --  @# pprof自定义server 转化为http server以满足transport.Server接口
	pprofServer := http.Server(*ps)
	// --  @# 追加http及grpc是否都启动的控制，都不启动时panic
	var ss []transport.Server

	if conf.Http.Enabled {
		log.NewHelper(logger).Infow("Kind", "Server", "Http", "enable")
		ss = append(ss, hs)
	}
	if conf.Grpc.Enabled {
		log.NewHelper(logger).Infow("Kind", "Server", "Grpc", "enable")
		ss = append(ss, gs)
	}
	if conf.Pprof.Enabled {
		log.NewHelper(logger).Infow("Kind", "Server", "Pprof", "enable")
		ss = append(ss, &pprofServer)
	}
	if 0 == len(ss) {
		log.NewHelper(logger).Errorw("Kind", "Server", "Error", "No Server Enabled!")
		panic("No Server Enabled! Please recheck server config!")
	}
	return kratos.New(
		kratos.ID(id),
		kratos.Name(conf.Name),
		kratos.Version(conf.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(ss...),
		kratos.Registrar(r),
	)
}

func main() {
	flag.Parse()

	// --  @# 启动时启动日志配置
	logger := boot.NewBootLogger(id).Run()

	// --  @# 启动时加载配置- consul kv
	c, bc := boot.NewBootConfig(address, schema, path).Run()
	defer func(c config.Config) {
		err := c.Close()
		if err != nil {
			panic(err)
		}
	}(c)

	// --  @# 启动时启动链路追踪配置
	boot.NewBootTrace(&bc).Run()

	// --  @# wire
	app, cleanup, err := wireApp(bc.Server, bc.Data, logger, bc.Registry)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
