package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	v1 "github.com/orzkratos/demokratos/demo2kratos/api/helloworld/v1"
	"github.com/orzkratos/demokratos/demo2kratos/internal/biz"
	"github.com/orzkratos/tracekratos"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc  *biz.GreeterUsecase
	log *log.Helper
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger) *GreeterService {
	return &GreeterService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	// Demo GetTraceID feature from tracekratos
	// 演示 tracekratos 的 GetTraceID 功能
	traceID := tracekratos.GetTraceID(ctx)
	s.log.WithContext(ctx).Infof("Processing request with trace ID: %s", traceID)

	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}
