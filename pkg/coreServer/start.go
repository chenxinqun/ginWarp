package coreServer

import (
	"errors"
	"flag"
	"fmt"
	"github.com/chenxinqun/ginWarpPkg/errno"
	"log"
	"net/http"
	"reflect"

	"github.com/chenxinqun/ginWarp/configs"
	"github.com/chenxinqun/ginWarp/pkg/core"

	"github.com/chenxinqun/ginWarpPkg/businessCodex"
	"github.com/chenxinqun/ginWarpPkg/datax/etcdx"
	"github.com/chenxinqun/ginWarpPkg/httpx/validation"
	logger "github.com/chenxinqun/ginWarpPkg/loggerx"
	"go.uber.org/zap"
)

type Args struct {
	env          string
	addr         string
	serviceName  string
	viewVersion  bool
	noHttpMod    bool
	closeFileLog bool
}

func (a Args) Env() string {
	return a.env
}

func (a Args) Addr() string {
	return a.addr
}

func (a Args) ServiceName() string {
	return a.serviceName
}

func (a Args) ViewVersion() bool {
	return a.viewVersion
}

func (a Args) NoHttpMod() bool {
	return a.noHttpMod
}

func (a Args) CloseFileLog() bool {
	return a.closeFileLog
}

var (
	args     = Args{}
	resource *Resource
)

func ReadArgs() Args {
	return args
}

func ReadResource() Resource {
	return *resource
}

func init() {
	flag.StringVar(&args.env, "env", "", "请输入运行环境如: \"-env dev\", "+
		"默认值为 \"dev\" \n dev:开发环境\n test:测试环境\n pre:预发布环境\n pro:正式环境\n")
	flag.StringVar(&args.addr, "addr", "", "请输入自定义的监听IP和端口, 如: \"-addr :8888\" \n "+
		", 如果要写IP则如:\"127.0.0.1:8888\" \n 服务启动一般会采用配置文件中的监听端口. 如果在命令行中输入此参数, 则会覆盖配配置文件中的\"Register.Listen\"配置项")
	flag.StringVar(&args.serviceName, "name", "", "请输入自定义的服务名称, 如: \"-name ginWarp-example\" \n "+
		", \n 服务启动一般会采用配置文件中的设置的服务名称.如果在命令行中输入此参数, 则会覆盖配配置文件中的\"Register.Name\"配置项")
	flag.BoolVar(&args.viewVersion, "version", false, "查询程序版本号")
	flag.BoolVar(&args.viewVersion, "nohttp", false, "不启动http接口, 当前程序以自定义的方式启动. 预留参数, 具体模式需要业务去实现.")
	flag.BoolVar(&args.closeFileLog, "closeFileLog", false, "不要文件日志, 只通过终端的方式输出日志.")

}

func FlagInit() Args {
	flag.Parse()
	return args
}

func StartHttpServer(r *Resource, noHttpTask func() error, routerHandlers ...core.RouterHandler) {
	if args.NoHttpMod() {
		log.Printf("当前程序以noHttp模式启动...")
		if noHttpTask == nil {
			log.Fatalf("noHttp模式, 请传入noHttpTask函数")
		}
	} else {
		log.Printf("当前程序以Http模式启动...")
	}
	if r == nil {
		r = new(Resource)
	}
	// 缓存配置源
	resource = r

	// 命令行参数优先, 传入的参数其次
	if args.Env() != "" {
		r.Env = args.Env()
	}
	// 命令行参数优先, 传入的参数其次
	if args.Addr() != "" {
		r.Addr = args.Addr()
	}

	// 命令行参数优先, 传入参数其次
	if args.ServiceName() != "" {
		r.ServiceName = args.ServiceName()
	}
	// 注册serviceName到框架
	core.SetServiceName(r.ServiceName)
	// 初始化服务启动必要数据
	InitServerManual(r.Env, r.ServiceName, r.Addr, r.BuildDate, r.BuildVersion, args.CloseFileLog())

	// 打印logo, 后期可以做一个产品的logo.
	fmt.Println(configs.UI)

	if reflect.DeepEqual(r.Cfg, configs.Config{}) {
		// 获取初始化之后的配置
		r.Cfg = configs.Default()
	}
	if args.ViewVersion() {
		log.Printf("当前编译文件版本号:%s, 编译时间:%s", r.BuildVersion, r.BuildDate)
		log.Fatalf("当前配置文件版本号:%s", r.Cfg.Register.ProjectVersion)
	}
	// 初始化 logger
	var (
		err     error
		loggers *zap.Logger
	)
	// 使用传进来的日志
	loggers = r.Logger
	if loggers == nil {
		// 如果没有传日志, 使用配置文件日志
		loggers, err = logger.New(&r.Cfg.Loggers)
		r.Logger = loggers
		if err != nil {
			panic(err)
		}
	}
	// 初始化etcd连接
	err = InitEtcdManual(r.Cfg, loggers)
	if err != nil {
		r.Logger.Fatal("初始化失败 ", zap.Error(err))
	}
	// 缓存etcdRepo
	r.EtcdRepo = etcdx.Default()
	// 从配置中心更新配置
	err = InitConfiguringManual(r.Cfg, r.EtcdRepo)
	if err != nil {
		r.Logger.Fatal("初始化失败 ", zap.Error(err))
	}
	// 再次初始化配置, 用默认配置, 补全配置中心没有设置过的必要配置项
	configs.InitConf(args.Addr(), args.ServiceName(), args.Env(), r.BuildDate, r.BuildVersion)
	// 获取更新后的配置
	r.Cfg = configs.Default()
	// 初始化验证器
	validation.InitValidation(r.Cfg.Language)
	businessCodex.SetLang(r.Cfg.Language)
	// 更新配置后, 重新创建日志
	logCfg := &r.Cfg.Loggers
	loggers, err = logger.New(logCfg)

	if err != nil {
		panic(err)
	}
	loggers.Info("日志创建")
	// 更新全局变量日志
	logger.SetDefault(loggers)
	if r.EtcdRepo != nil {
		// 替换日志对象为重新创建后的日志.
		r.EtcdRepo.GetRepo().Logger = loggers
	}
	r.Logger = loggers

	// 刷新日志缓存到磁盘
	defer func() {
		_ = r.Logger.Sync()
	}()
	// 打印一些信息
	log.Printf("* [register microservice model \"%v\"]", r.Cfg.MicroserviceModel)
	log.Printf("* [register business platform \"%s\"]", r.Cfg.BusinessPlatform)
	log.Printf("* [register service \"%s\"]", r.Cfg.Register.Name)
	log.Printf("* [register listen \"%s\"]", r.Cfg.Register.Listen)
	log.Printf("* [register environment \"%s\"]", r.Cfg.Register.Namespace)

	// 前置初始化
	InitHandlerExec(r.Cfg, r.Logger, initBeforeHandlers...)
	//用户自定的初始化
	InitHandlerExec(r.Cfg, r.Logger, initHandlers...)

	// 非http模式
	if args.NoHttpMod() {
		err = noHttpTask()
		if err != nil {
			log.Fatalf("非Http模式程序运行时报错退出:%v", errno.WithStack(err))
		}
		return
	} else {

		// 初始化 HTTP 服务器
		s, err := NewHTTPServer(r, routerHandlers...)
		if err != nil {
			r.Logger.Fatal("初始化失败", zap.Error(err))
		}

		// 初始化完毕, 缓存server
		SetServer(s)

		// 创建http coreServer
		httpServer := &http.Server{
			Addr:    r.Cfg.Register.Listen,
			Handler: s.Mux,
		}
		// 缓存http coreServer
		SetHttpServer(httpServer)

		// 启动 http coreServer
		go func() {
			if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				r.Logger.Fatal("http coreServer startup err", zap.Error(err))
			}
		}()

		/* 后置初始化 在 http coreServer 启动之后 主要是服务注册和发现 */
		InitHandlerExec(r.Cfg, r.Logger, initAfterHandlers...)
		log.Printf("* [Host %s] \n", r.Cfg.Register.Url())
	}

}
