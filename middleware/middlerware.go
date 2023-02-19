package middleware

import (
	"github.com/chenxinqun/ginWarpPkg/datax/mysqlx"
	"github.com/chenxinqun/ginWarpPkg/datax/redisx"
	"go.uber.org/zap"
)

type Middleware struct {
	Logger *zap.Logger
	Mysql  mysqlx.Repo
	Redis  redisx.Repo
}

func New(logger *zap.Logger, redis redisx.Repo, repo mysqlx.Repo) *Middleware {
	return &Middleware{
		Logger: logger,
		Mysql:  mysql,
		Redis:  redis,
	}
}
