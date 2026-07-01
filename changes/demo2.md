# Changes

Code differences compared to source project.

## cmd/demo2kratos/main.go (+13 -18)

```diff
@@ -2,20 +2,19 @@
 
 import (
 	"flag"
-	"log/slog"
 	"os"
 
-	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
 	"github.com/go-kratos/kratos/v3"
 	"github.com/go-kratos/kratos/v3/config"
 	"github.com/go-kratos/kratos/v3/config/file"
-	"github.com/go-kratos/kratos/v3/log"
 	"github.com/go-kratos/kratos/v3/transport/grpc"
 	"github.com/go-kratos/kratos/v3/transport/http"
 	"github.com/yylego/done"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 	"github.com/yylego/must"
 	"github.com/yylego/rese"
+	"github.com/yylego/zaplog"
 
 	_ "go.uber.org/automaxprocs"
 )
@@ -34,13 +33,13 @@
 	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
 }
 
-func newApp(logger *slog.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
+func newApp(gs *grpc.Server, hs *http.Server, zap匝普日志 *zapzhkratos.T匝普日志) *kratos.App {
 	return kratos.New(
 		kratos.ID(done.VCE(os.Hostname()).Omit()),
 		kratos.Name(Name),
 		kratos.Version(Version),
 		kratos.Metadata(map[string]string{}),
-		kratos.Logger(logger),
+		kratos.Logger(zap匝普日志.Get奎沱日志("网络服务")),
 		kratos.Server(
 			gs,
 			hs,
@@ -50,18 +49,14 @@
 
 func main() {
 	flag.Parse()
-	logger := log.NewLogger(
-		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
-			AddSource: true,
-			Level:     slog.LevelInfo,
-		}),
-		log.WithExtractor(tracing.TraceAttrs),
-	).With(
-		slog.String("service.id", done.VCE(os.Hostname()).Omit()),
-		slog.String("service.name", Name),
-		slog.String("service.version", Version),
-	)
-	log.SetDefault(logger)
+
+	// demo2 uses Get奎沱秘书 to derive a *slog.Logger for the startup logs
+	// demo2 使用 Get奎沱秘书 派生 *slog.Logger 打印启动日志
+	zapKratos := zapzhkratos.New匝普日志(zaplog.LOGGER, zapzhkratos.New日志配置())
+	slog := zapKratos.Get奎沱秘书("启动日志")
+	slog.Info("服务版本", "version", Version)
+	slog.Info("配置路径", "config", flagconf)
+
 	c := config.New(
 		config.WithSource(
 			file.NewSource(flagconf),
@@ -74,7 +69,7 @@
 	var cfg conf.Bootstrap
 	must.Done(c.Scan(&cfg))
 
-	app, cleanup := rese.V2(wireApp(cfg.Server, cfg.Data, logger))
+	app, cleanup := rese.V2(wireApp(cfg.Server, cfg.Data, zapKratos))
 	defer cleanup()
 
 	// start and wait for stop signal
```

## cmd/demo2kratos/wire.go (+2 -3)

```diff
@@ -5,8 +5,6 @@
 package main
 
 import (
-	"log/slog"
-
 	"github.com/go-kratos/kratos/v3"
 	"github.com/google/wire"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/biz"
@@ -14,9 +12,10 @@
 	"github.com/yylego/kratos-examples/demo2kratos/internal/data"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/server"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/service"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 )
 
 // wireApp init kratos application.
-func wireApp(*conf.Server, *conf.Data, *slog.Logger) (*kratos.App, func(), error) {
+func wireApp(*conf.Server, *conf.Data, *zapzhkratos.T匝普日志) (*kratos.App, func(), error) {
 	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
 }
```

## cmd/demo2kratos/wire_gen.go (+8 -8)

```diff
@@ -13,7 +13,7 @@
 	"github.com/yylego/kratos-examples/demo2kratos/internal/data"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/server"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/service"
-	"log/slog"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 )
 
 import (
@@ -23,20 +23,20 @@
 // Injectors from wire.go:
 
 // wireApp init kratos application.
-func wireApp(confServer *conf.Server, confData *conf.Data, logger *slog.Logger) (*kratos.App, func(), error) {
-	dataData, cleanup, err := data.NewData(confData, logger)
+func wireApp(confServer *conf.Server, confData *conf.Data, t匝普日志 *zapzhkratos.T匝普日志) (*kratos.App, func(), error) {
+	dataData, cleanup, err := data.NewData(confData, t匝普日志)
 	if err != nil {
 		return nil, nil, err
 	}
-	articleUsecase, err := biz.NewArticleUsecase(dataData, logger)
+	articleUsecase, err := biz.NewArticleUsecase(dataData, t匝普日志)
 	if err != nil {
 		cleanup()
 		return nil, nil, err
 	}
-	articleService := service.NewArticleService(articleUsecase)
-	grpcServer := server.NewGRPCServer(confServer, articleService, logger)
-	httpServer := server.NewHTTPServer(confServer, articleService, logger)
-	app := newApp(logger, grpcServer, httpServer)
+	articleService := service.NewArticleService(articleUsecase, t匝普日志)
+	grpcServer := server.NewGRPCServer(confServer, articleService, t匝普日志)
+	httpServer := server.NewHTTPServer(confServer, articleService, t匝普日志)
+	app := newApp(grpcServer, httpServer, t匝普日志)
 	return app, func() {
 		cleanup()
 	}, nil
```

## internal/biz/article.go (+5 -4)

```diff
@@ -8,6 +8,7 @@
 	"github.com/yylego/kratos-ebz/ebzkratos"
 	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/data"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 	"github.com/yylego/must"
 	"gorm.io/gorm"
 	"gorm.io/gorm/clause"
@@ -33,14 +34,14 @@
 	slog *slog.Logger
 }
 
-func NewArticleUsecase(data *data.Data, logger *slog.Logger) (*ArticleUsecase, error) {
+func NewArticleUsecase(data *data.Data, zap匝普日志 *zapzhkratos.T匝普日志) (*ArticleUsecase, error) {
 	// Migrate the owned table plus the mirrored students table (needed in the
 	// existence check); both services share one database
 	// 建好本服务拥有的 articles 表，外加镜像的 students 表（供存在性校验用）
 	if err := data.DB().AutoMigrate(&Article{}, &Student{}); err != nil {
 		return nil, err
 	}
-	return &ArticleUsecase{data: data, slog: logger}, nil
+	return &ArticleUsecase{data: data, slog: zap匝普日志.Get奎沱秘书("业务逻辑")}, nil
 }
 
 func (uc *ArticleUsecase) CreateArticle(ctx context.Context, a *Article) (*Article, *ebzkratos.Ebz) {
@@ -68,7 +69,7 @@
 		}
 		return nil, ebzkratos.New(pb.ErrorArticleCreateFailure("create article: %v", err))
 	}
-	uc.slog.InfoContext(ctx, "created article", "id", res.ID, "student_id", res.StudentID)
+	uc.slog.InfoContext(ctx, "已创建文章", "id", res.ID, "student_id", res.StudentID)
 	return res, nil
 }
 
@@ -127,7 +128,7 @@
 	if del.RowsAffected == 0 {
 		return ebzkratos.New(pb.ErrorArticleNotFound("article %d not found", id))
 	}
-	uc.slog.InfoContext(ctx, "deleted article", "id", id)
+	uc.slog.InfoContext(ctx, "已删除文章", "id", id)
 	return nil
 }
 
```

## internal/data/data.go (+5 -4)

```diff
@@ -1,10 +1,9 @@
 package data
 
 import (
-	"log/slog"
-
 	"github.com/google/wire"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 	"github.com/yylego/must"
 	"github.com/yylego/rese"
 	"gorm.io/driver/postgres"
@@ -24,11 +23,13 @@
 	return d.db
 }
 
-func NewData(c *conf.Data, logger *slog.Logger) (*Data, func(), error) {
+func NewData(c *conf.Data, zap匝普日志 *zapzhkratos.T匝普日志) (*Data, func(), error) {
+	slog := zap匝普日志.Get奎沱秘书("数据层")
+	slog.Debug("准备链接数据资源")
 	must.Same(c.Database.Driver, "postgres")
 	db := rese.P1(gorm.Open(postgres.Open(c.Database.Source), &gorm.Config{}))
 	cleanup := func() {
-		logger.Info("closing the data resources")
+		slog.Info("准备关闭数据资源")
 		_ = rese.P1(db.DB()).Close()
 	}
 	return &Data{db: db}, cleanup, nil
```

## internal/server/grpc.go (+4 -3)

```diff
@@ -1,19 +1,20 @@
 package server
 
 import (
-	"log/slog"
-
+	"github.com/go-kratos/kratos/v3/middleware/logging"
 	"github.com/go-kratos/kratos/v3/middleware/recovery"
 	"github.com/go-kratos/kratos/v3/transport/grpc"
 	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/service"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 )
 
-func NewGRPCServer(c *conf.Server, article *service.ArticleService, logger *slog.Logger) *grpc.Server {
+func NewGRPCServer(c *conf.Server, article *service.ArticleService, zap匝普日志 *zapzhkratos.T匝普日志) *grpc.Server {
 	var opts = []grpc.ServerOption{
 		grpc.Middleware(
 			recovery.Recovery(),
+			logging.Server(zap匝普日志.Get奎沱日志("GRPC请求日志")),
 		),
 	}
 	if c.Grpc.Network != "" {
```

## internal/server/http.go (+4 -3)

```diff
@@ -1,19 +1,20 @@
 package server
 
 import (
-	"log/slog"
-
+	"github.com/go-kratos/kratos/v3/middleware/logging"
 	"github.com/go-kratos/kratos/v3/middleware/recovery"
 	"github.com/go-kratos/kratos/v3/transport/http"
 	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/service"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 )
 
-func NewHTTPServer(c *conf.Server, article *service.ArticleService, logger *slog.Logger) *http.Server {
+func NewHTTPServer(c *conf.Server, article *service.ArticleService, zap匝普日志 *zapzhkratos.T匝普日志) *http.Server {
 	var opts = []http.ServerOption{
 		http.Middleware(
 			recovery.Recovery(),
+			logging.Server(zap匝普日志.Get奎沱日志("HTTP请求日志")),
 		),
 	}
 	if c.Http.Network != "" {
```

## internal/service/article.go (+11 -3)

```diff
@@ -2,22 +2,29 @@
 
 import (
 	"context"
+	"log/slog"
 
 	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/biz"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 )
 
 type ArticleService struct {
 	pb.UnimplementedArticleServiceServer
 
-	uc *biz.ArticleUsecase
+	uc   *biz.ArticleUsecase
+	slog *slog.Logger
 }
 
-func NewArticleService(uc *biz.ArticleUsecase) *ArticleService {
-	return &ArticleService{uc: uc}
+func NewArticleService(uc *biz.ArticleUsecase, zap匝普日志 *zapzhkratos.T匝普日志) *ArticleService {
+	return &ArticleService{
+		uc:   uc,
+		slog: zap匝普日志.Get奎沱秘书("服务层"),
+	}
 }
 
 func (s *ArticleService) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleReply, error) {
+	s.slog.InfoContext(ctx, "收到请求: create-article")
 	if req.Title == "" {
 		return nil, pb.ErrorBadParam("TITLE IS REQUIRED")
 	}
@@ -32,6 +39,7 @@
 	if ebz != nil {
 		return nil, ebz.Erk
 	}
+	s.slog.InfoContext(ctx, "返回响应: create-article")
 	return &pb.CreateArticleReply{Article: &pb.ArticleInfo{Id: v.ID, Title: v.Title, Content: v.Content, StudentId: v.StudentID}}, nil
 }
 
```

