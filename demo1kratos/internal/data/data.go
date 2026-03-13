package data

import (
	"github.com/google/wire"
	"github.com/yylego/kratos-examples/demo1kratos/internal/conf"
	"github.com/yylego/kratos-zapzh/zapzhkratos"
	"github.com/yylego/must"
	"github.com/yylego/rese"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(NewData)

type Data struct {
	db *gorm.DB
}

func NewData(c *conf.Data, zap匝普日志 *zapzhkratos.T匝普日志) (*Data, func(), error) {
	zapLog := zap匝普日志.Sub模块匝普()
	zapLog.SUG.Debugln("准备链接数据资源")
	must.Same(c.Database.Driver, "sqlite3")
	db := rese.P1(gorm.Open(sqlite.Open(c.Database.Source), &gorm.Config{}))
	cleanup := func() {
		zapLog.SUG.Info("准备关闭数据资源")
		_ = rese.P1(db.DB()).Close()
	}
	return &Data{db: db}, cleanup, nil
}
