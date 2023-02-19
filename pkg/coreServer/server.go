package coreServer

import (
	"github.com/chenxinqun/ginWarpPkg/datax/mysqlx"
	"github.com/chenxinqun/ginWarpPkg/datax/redisx"
	"github.com/chenxinqun/ginWarpPkg/errno"

	"github.com/chenxinqun/ginWarp/configs"
	"github.com/chenxinqun/ginWarp/example/router"
	"github.com/chenxinqun/ginWarp/pkg/core"

	"github.com/chenxinqun/ginWarpPkg/datax/etcdx"
	"github.com/chenxinqun/ginWarpPkg/datax/kafkax"
	"github.com/chenxinqun/ginWarpPkg/datax/mongox"
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
	logger "github.com/chenxinqun/ginWarpPkg/loggerx"
	"go.uber.org/zap"
)

type Server struct {
	Mux               mux.IMux
	EtcdRepo          etcdx.Repo
	MysqlRepo         mysqlx.Repo
	RedisRepo         redisx.Repo
	MongoRepo         mongox.Repo
	ProducerRepo      kafkax.ProducerRepo
	AsyncProducerRepo kafkax.AsyncProducerRepo
	ConsumerGroupRepo kafkax.ConsumerGroupRepo
	Logger            *zap.Logger
	Cfg               configs.Config
}

type Resource struct {
	Env          string
	ServiceName  string
	Addr         string
	BuildDate    string // 编译时,由CI/CD工具复写此变量, 并在main函数中传进来
	BuildVersion string // 编译时,由CI/CD工具复写此变量, 并在main函数中传进来
	// 一些 core 中的组件的开关
	MuxOptions []mux.OptionHandler
	// 配置
	Cfg      configs.Config
	EtcdRepo etcdx.Repo
	// 日志
	Logger *zap.Logger
}

func NewHTTPServer(r *Resource, routerHandlers ...core.RouterHandler) (*Server, error) {
	var server *Server
	if r.Logger == nil {
		r.Logger = logger.Default()
	}
	server = new(Server)
	routerR := new(core.RouterResource)
	// 初始化 Mux
	muxR := &mux.Resource{
		Env:           r.Cfg.Env,
		Logger:        r.Logger,
		ProjectListen: r.Cfg.Register.Listen,
	}
	m, err := core.NewMux(muxR, r.MuxOptions...)

	if err != nil {
		return nil, errno.Errorf("http coreServer init err: %v", err)
	}

	routerR.Mux = m

	if len(routerHandlers) > 0 {
		// 设置 API 路由
		for _, handler := range routerHandlers {
			handler(routerR)
		}
	} else {
		// 开启默认的 example 路由
		router.SetApiRouter(routerR)
	}

	server.Mux = m
	server.EtcdRepo = etcdx.Default()
	server.MysqlRepo = mysqlx.Default()
	server.RedisRepo = redisx.Default()
	server.MongoRepo = mongox.Default()
	server.ProducerRepo = kafkax.DefaultProducer()
	server.AsyncProducerRepo = kafkax.DefaultAsyncProducer()
	server.ConsumerGroupRepo = kafkax.DefaultConsumerGroup()
	server.Logger = r.Logger
	server.Cfg = r.Cfg
	return server, nil
}
