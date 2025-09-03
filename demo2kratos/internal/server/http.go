package server

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/uuid"
	v1 "github.com/orzkratos/demokratos/demo2kratos/api/helloworld/v1"
	"github.com/orzkratos/demokratos/demo2kratos/internal/conf"
	"github.com/orzkratos/demokratos/demo2kratos/internal/service"
	"github.com/orzkratos/tracekratos"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			NewTraceMiddleware(logger), //在请求逻辑执行前打印日志，显示请求参数和追踪信息
			logging.Server(logger),     //在请求逻辑执行后打印日志，显示执行结果的错误码和状态码
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

func NewTraceMiddleware(logger log.Logger) middleware.Middleware {
	config := tracekratos.NewConfig("TRACE_ID")
	config.NewTraceID = func(ctx context.Context) string {
		return "TRACE-ID-" + strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + uuid.New().String() + "-BBB"
	}
	return tracekratos.NewTraceMiddleware(config, logger)
}
