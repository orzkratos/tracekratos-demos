# Changes

Code differences compared to source project demokratos.

## cmd/demo2kratos/wire_gen.go (+1 -1)

```diff
@@ -25,7 +25,7 @@
 	}
 	greeterRepo := data.NewGreeterRepo(dataData, logger)
 	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
-	greeterService := service.NewGreeterService(greeterUsecase)
+	greeterService := service.NewGreeterService(greeterUsecase, logger)
 	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
 	httpServer := server.NewHTTPServer(confServer, greeterService, logger)
 	app := newApp(logger, grpcServer, httpServer)
```

## internal/server/http.go (+26 -0)

```diff
@@ -1,12 +1,20 @@
 package server
 
 import (
+	"context"
+	"strconv"
+	"time"
+
 	"github.com/go-kratos/kratos/v2/log"
+	"github.com/go-kratos/kratos/v2/middleware"
+	"github.com/go-kratos/kratos/v2/middleware/logging"
 	"github.com/go-kratos/kratos/v2/middleware/recovery"
 	"github.com/go-kratos/kratos/v2/transport/http"
+	"github.com/google/uuid"
 	v1 "github.com/orzkratos/demokratos/demo2kratos/api/helloworld/v1"
 	"github.com/orzkratos/demokratos/demo2kratos/internal/conf"
 	"github.com/orzkratos/demokratos/demo2kratos/internal/service"
+	"github.com/orzkratos/tracekratos"
 )
 
 // NewHTTPServer new an HTTP server.
@@ -14,6 +22,8 @@
 	var opts = []http.ServerOption{
 		http.Middleware(
 			recovery.Recovery(),
+			NewTraceMiddleware(logger), //在请求逻辑执行前打印日志，显示请求参数和追踪信息
+			logging.Server(logger),     //在请求逻辑执行后打印日志，显示执行结果的错误码和状态码
 		),
 	}
 	if c.Http.Network != "" {
@@ -28,4 +38,20 @@
 	srv := http.NewServer(opts...)
 	v1.RegisterGreeterHTTPServer(srv, greeter)
 	return srv
+}
+
+func NewTraceMiddleware(logger log.Logger) middleware.Middleware {
+	// Demo tracekratos features using function options
+	// 演示 tracekratos 的功能选项
+	config := tracekratos.NewConfig("TRACE_ID",
+		tracekratos.WithLogLevel(log.LevelInfo),
+		tracekratos.WithLogReply(true),
+		tracekratos.WithNewTraceID(func(ctx context.Context) string {
+			return "TRACE-ID-" + strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + uuid.New().String() + "-BBB"
+		}),
+		tracekratos.WithFormatArgs(func(req any) string {
+			return tracekratos.ExtractArgs(req)
+		}),
+	)
+	return tracekratos.NewTraceMiddleware(config, logger)
 }
```

## internal/service/greeter.go (+14 -3)

```diff
@@ -3,24 +3,35 @@
 import (
 	"context"
 
+	"github.com/go-kratos/kratos/v2/log"
 	v1 "github.com/orzkratos/demokratos/demo2kratos/api/helloworld/v1"
 	"github.com/orzkratos/demokratos/demo2kratos/internal/biz"
+	"github.com/orzkratos/tracekratos"
 )
 
 // GreeterService is a greeter service.
 type GreeterService struct {
 	v1.UnimplementedGreeterServer
 
-	uc *biz.GreeterUsecase
+	uc  *biz.GreeterUsecase
+	log *log.Helper
 }
 
 // NewGreeterService new a greeter service.
-func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
-	return &GreeterService{uc: uc}
+func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger) *GreeterService {
+	return &GreeterService{
+		uc:  uc,
+		log: log.NewHelper(logger),
+	}
 }
 
 // SayHello implements helloworld.GreeterServer.
 func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
+	// Demo GetTraceID feature from tracekratos
+	// 演示 tracekratos 的 GetTraceID 功能
+	traceID := tracekratos.GetTraceID(ctx)
+	s.log.WithContext(ctx).Infof("Processing request with trace ID: %s", traceID)
+
 	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
 	if err != nil {
 		return nil, err
```

