package traceinfo_test

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/orzkratos/demokratos/demo2kratos/internal/pkg/traceinfo"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// 其实也没必要对 middleware 进行单元测试的，这里只是自己尝试写写，掌握点新技能
func TestNewTraceMiddleware(t *testing.T) {
	var rawFunc middleware.Handler = func(ctx context.Context, req any) (any, error) {
		return structpb.NewStruct(map[string]any{
			"A": 10,
			"B": 11,
			"C": 12,
		})
	}

	runFunc := traceinfo.NewTraceMiddleware(log.DefaultLogger)(rawFunc)

	ctx := context.Background()
	res, err := runFunc(ctx, &emptypb.Empty{})
	require.NoError(t, err)
	t.Log(neatjsons.S(res))

	object, ok := res.(*structpb.Struct)
	require.True(t, ok)
	require.Equal(t, map[string]any{
		"A": 10.0,
		"B": 11.0,
		"C": 12.0,
	}, object.AsMap())
}
