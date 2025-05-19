package traceinfo

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/google/uuid"
	"github.com/orzkratos/tracekratos"
)

func NewTraceMiddleware(logger log.Logger) middleware.Middleware {
	config := tracekratos.NewConfig("TRACE_ID")
	config.NewTraceID = func(ctx context.Context) string {
		return "TRACE-ID-" + strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + strings.ToUpper(uuid.New().String()) + "-AAA"
	}
	return tracekratos.NewTraceMiddleware(config, logger)
}
