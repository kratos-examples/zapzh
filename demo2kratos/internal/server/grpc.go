package server

import (
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
	"github.com/yylego/kratos-examples/demo2kratos/internal/service"
	"github.com/yylego/kratos-zapzh/zapzhkratos"
)

func NewGRPCServer(c *conf.Server, article *service.ArticleService, zap匝普日志 *zapzhkratos.T匝普日志) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(zap匝普日志.Get奎沱日志("GRPC请求")),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Address != "" {
		opts = append(opts, grpc.Address(c.Grpc.Address))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	pb.RegisterArticleServiceServer(srv, article)
	return srv
}
