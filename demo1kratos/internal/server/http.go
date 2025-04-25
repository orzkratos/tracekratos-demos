package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	v1 "github.com/orzkratos/demokratos/demo1kratos/api/helloworld/v1"
	"github.com/orzkratos/demokratos/demo1kratos/internal/conf"
	"github.com/orzkratos/demokratos/demo1kratos/internal/service"
	"github.com/orzkratos/tracekratos"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracekratos.NewTraceMiddleware(tracekratos.NewConfig("TRACE_ID"), logger), //在请求逻辑执行前打印日志，显示请求参数和追踪信息
			logging.Server(logger), //在请求逻辑执行后打印日志，显示执行结果的错误码和状态码
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	return srv
}
