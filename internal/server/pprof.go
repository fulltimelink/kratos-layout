package server

import (
	"github.com/fulltimelink/kratos-layout/internal/conf"
	"net/http/pprof"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// PProfServer --  @# 定义新的http.Server类型，规避wire只能依赖注入一个Http Server，wire中将其再转为http.Server
type PProfServer http.Server

func RegisterPprof(s *http.Server) {
	s.HandleFunc("/debug/pprof", pprof.Index)
	s.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	s.HandleFunc("/debug/pprof/profile", pprof.Profile)
	s.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	s.HandleFunc("/debug/pprof/trace", pprof.Trace)
	s.HandleFunc("/debug/allocs", pprof.Handler("allocs").ServeHTTP)
	s.HandleFunc("/debug/block", pprof.Handler("block").ServeHTTP)
	s.HandleFunc("/debug/goroutine", pprof.Handler("goroutine").ServeHTTP)
	s.HandleFunc("/debug/heap", pprof.Handler("heap").ServeHTTP)
	s.HandleFunc("/debug/mutex", pprof.Handler("mutex").ServeHTTP)
	s.HandleFunc("/debug/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}

func NewPProfServer(conf *conf.Server) *PProfServer {
	httpSrv := http.NewServer(
		http.Address(conf.Pprof.Addr),
	)
	RegisterPprof(httpSrv)
	pprofSrv := PProfServer(*httpSrv)

	return &pprofSrv
}
