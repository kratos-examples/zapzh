# Changes

Code differences compared to source project.

## cmd/demo2kratos/main.go (+14 -14)

```diff
@@ -7,14 +7,15 @@
 	"github.com/go-kratos/kratos/v2"
 	"github.com/go-kratos/kratos/v2/config"
 	"github.com/go-kratos/kratos/v2/config/file"
-	"github.com/go-kratos/kratos/v2/log"
-	"github.com/go-kratos/kratos/v2/middleware/tracing"
 	"github.com/go-kratos/kratos/v2/transport/grpc"
 	"github.com/go-kratos/kratos/v2/transport/http"
 	"github.com/yylego/done"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 	"github.com/yylego/must"
 	"github.com/yylego/rese"
+	"github.com/yylego/zaplog"
+	_ "go.uber.org/automaxprocs"
 )
 
 // go build -ldflags "-X main.Version=x.y.z"
@@ -31,13 +32,13 @@
 	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
 }
 
-func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
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
@@ -47,15 +48,14 @@
 
 func main() {
 	flag.Parse()
-	logger := log.With(log.NewStdLogger(os.Stdout),
-		"ts", log.DefaultTimestamp,
-		"caller", log.DefaultCaller,
-		"service.id", kratos.ID(done.VCE(os.Hostname()).Omit()),
-		"service.name", Name,
-		"service.version", Version,
-		"trace.id", tracing.TraceID(),
-		"span.id", tracing.SpanID(),
-	)
+
+	// demo2 uses Get奎沱秘书 to get *log.Helper (Kratos style)
+	// demo2 使用 Get奎沱秘书 获取 *log.Helper（Kratos 风格）
+	zapKratos := zapzhkratos.New匝普日志(zaplog.LOGGER, zapzhkratos.New日志配置())
+	slog := zapKratos.Get奎沱秘书("启动日志")
+	slog.Infof("服务版本: %s", Version)
+	slog.Infof("配置路径: %s", flagconf)
+
 	c := config.New(
 		config.WithSource(
 			file.NewSource(flagconf),
@@ -68,7 +68,7 @@
 	var cfg conf.Bootstrap
 	must.Done(c.Scan(&cfg))
 
-	app, cleanup := rese.V2(wireApp(cfg.Server, cfg.Data, logger))
+	app, cleanup := rese.V2(wireApp(cfg.Server, cfg.Data, zapKratos))
 	defer cleanup()
 
 	// start and wait for stop signal
```

## cmd/demo2kratos/wire.go (+2 -2)

```diff
@@ -6,16 +6,16 @@
 
 import (
 	"github.com/go-kratos/kratos/v2"
-	"github.com/go-kratos/kratos/v2/log"
 	"github.com/google/wire"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/biz"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/data"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/server"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/service"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 )
 
 // wireApp init kratos application.
-func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
+func wireApp(*conf.Server, *conf.Data, *zapzhkratos.T匝普日志) (*kratos.App, func(), error) {
 	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
 }
```

## cmd/demo2kratos/wire_gen.go (+12 -8)

```diff
@@ -7,27 +7,31 @@
 
 import (
 	"github.com/go-kratos/kratos/v2"
-	"github.com/go-kratos/kratos/v2/log"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/biz"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/data"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/server"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/service"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 )
 
+import (
+	_ "go.uber.org/automaxprocs"
+)
+
 // Injectors from wire.go:
 
 // wireApp init kratos application.
-func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
-	dataData, cleanup, err := data.NewData(confData, logger)
+func wireApp(confServer *conf.Server, confData *conf.Data, t匝普日志 *zapzhkratos.T匝普日志) (*kratos.App, func(), error) {
+	dataData, cleanup, err := data.NewData(confData, t匝普日志)
 	if err != nil {
 		return nil, nil, err
 	}
-	articleUsecase := biz.NewArticleUsecase(dataData, logger)
-	articleService := service.NewArticleService(articleUsecase)
-	grpcServer := server.NewGRPCServer(confServer, articleService, logger)
-	httpServer := server.NewHTTPServer(confServer, articleService, logger)
-	app := newApp(logger, grpcServer, httpServer)
+	articleUsecase := biz.NewArticleUsecase(dataData, t匝普日志)
+	articleService := service.NewArticleService(articleUsecase, t匝普日志)
+	grpcServer := server.NewGRPCServer(confServer, articleService, t匝普日志)
+	httpServer := server.NewHTTPServer(confServer, articleService, t匝普日志)
+	app := newApp(grpcServer, httpServer, t匝普日志)
 	return app, func() {
 		cleanup()
 	}, nil
```

## internal/biz/article.go (+4 -3)

```diff
@@ -8,6 +8,7 @@
 	"github.com/yylego/kratos-ebz/ebzkratos"
 	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/data"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 )
 
 type Article struct {
@@ -19,11 +20,11 @@
 
 type ArticleUsecase struct {
 	data *data.Data
-	log  *log.Helper
+	slog *log.Helper
 }
 
-func NewArticleUsecase(data *data.Data, logger log.Logger) *ArticleUsecase {
-	return &ArticleUsecase{data: data, log: log.NewHelper(logger)}
+func NewArticleUsecase(data *data.Data, zap匝普日志 *zapzhkratos.T匝普日志) *ArticleUsecase {
+	return &ArticleUsecase{data: data, slog: zap匝普日志.Get奎沱秘书("业务逻辑")}
 }
 
 func (uc *ArticleUsecase) CreateArticle(ctx context.Context, a *Article) (*Article, *ebzkratos.Ebz) {
```

## internal/data/data.go (+5 -3)

```diff
@@ -1,9 +1,9 @@
 package data
 
 import (
-	"github.com/go-kratos/kratos/v2/log"
 	"github.com/google/wire"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 	"github.com/yylego/must"
 	"github.com/yylego/rese"
 	"gorm.io/driver/sqlite"
@@ -16,11 +16,13 @@
 	db *gorm.DB
 }
 
-func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
+func NewData(c *conf.Data, zap匝普日志 *zapzhkratos.T匝普日志) (*Data, func(), error) {
+	slog := zap匝普日志.Get奎沱秘书("数据层")
+	slog.Debug("准备链接数据资源")
 	must.Same(c.Database.Driver, "sqlite3")
 	db := rese.P1(gorm.Open(sqlite.Open(c.Database.Source), &gorm.Config{}))
 	cleanup := func() {
-		log.NewHelper(logger).Info("closing the data resources")
+		slog.Info("准备关闭数据资源")
 		_ = rese.P1(db.DB()).Close()
 	}
 	return &Data{db: db}, cleanup, nil
```

## internal/server/grpc.go (+4 -2)

```diff
@@ -1,18 +1,20 @@
 package server
 
 import (
-	"github.com/go-kratos/kratos/v2/log"
+	"github.com/go-kratos/kratos/v2/middleware/logging"
 	"github.com/go-kratos/kratos/v2/middleware/recovery"
 	"github.com/go-kratos/kratos/v2/transport/grpc"
 	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/service"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 )
 
-func NewGRPCServer(c *conf.Server, article *service.ArticleService, logger log.Logger) *grpc.Server {
+func NewGRPCServer(c *conf.Server, article *service.ArticleService, zap匝普日志 *zapzhkratos.T匝普日志) *grpc.Server {
 	var opts = []grpc.ServerOption{
 		grpc.Middleware(
 			recovery.Recovery(),
+			logging.Server(zap匝普日志.Get奎沱日志("GRPC请求")),
 		),
 	}
 	if c.Grpc.Network != "" {
```

## internal/server/http.go (+4 -2)

```diff
@@ -1,18 +1,20 @@
 package server
 
 import (
-	"github.com/go-kratos/kratos/v2/log"
+	"github.com/go-kratos/kratos/v2/middleware/logging"
 	"github.com/go-kratos/kratos/v2/middleware/recovery"
 	"github.com/go-kratos/kratos/v2/transport/http"
 	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/conf"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/service"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 )
 
-func NewHTTPServer(c *conf.Server, article *service.ArticleService, logger log.Logger) *http.Server {
+func NewHTTPServer(c *conf.Server, article *service.ArticleService, zap匝普日志 *zapzhkratos.T匝普日志) *http.Server {
 	var opts = []http.ServerOption{
 		http.Middleware(
 			recovery.Recovery(),
+			logging.Server(zap匝普日志.Get奎沱日志("HTTP请求")),
 		),
 	}
 	if c.Http.Network != "" {
```

## internal/service/article.go (+11 -3)

```diff
@@ -3,25 +3,33 @@
 import (
 	"context"
 
+	"github.com/go-kratos/kratos/v2/log"
 	pb "github.com/yylego/kratos-examples/demo2kratos/api/article"
 	"github.com/yylego/kratos-examples/demo2kratos/internal/biz"
+	"github.com/yylego/kratos-zapzh/zapzhkratos"
 )
 
 type ArticleService struct {
 	pb.UnimplementedArticleServiceServer
 
-	uc *biz.ArticleUsecase
+	uc   *biz.ArticleUsecase
+	slog *log.Helper
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
+	s.slog.WithContext(ctx).Infof("收到请求: create-article")
 	v, ebz := s.uc.CreateArticle(ctx, nil)
 	if ebz != nil {
 		return nil, ebz.Erk
 	}
+	s.slog.WithContext(ctx).Infof("返回响应: create-article")
 	return &pb.CreateArticleReply{Article: &pb.ArticleInfo{Id: v.ID, Title: v.Title, Content: v.Content, StudentId: v.StudentID}}, nil
 }
 
```

