package core

import (
	"fmt"
	"github.com/chenxinqun/ginWarpPkg/timex"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/chenxinqun/ginWarp/configs"

	"github.com/chenxinqun/ginWarpPkg/businessCodex"
	"github.com/chenxinqun/ginWarpPkg/errno"
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
	"github.com/chenxinqun/ginWarpPkg/sysx/environment"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var (
	defaultMux  mux.IMux
	serviceName string
)

type RouterHandler func(r *RouterResource)

type RouterResource struct {
	Mux mux.IMux
}

func DefaultMux() mux.IMux {
	return defaultMux
}
func SetServiceName(name string) {
	serviceName = name
}

func GetServiceName() string {
	return serviceName
}

func GetServiceNameUrl() string {
	return fmt.Sprintf("/%s", GetServiceName())
}

func NewMux(r *mux.Resource, options ...mux.OptionHandler) (mux.IMux, error) {
	logger := r.Logger
	if logger == nil {
		return nil, errno.NewError("logger required")
	}

	if environment.Active().IsDev() {
		// dev环境设置为debug模式
		gin.SetMode(gin.DebugMode)

	} else if environment.Active().IsTest() {
		// test环境设置为测试模式
		gin.SetMode(gin.TestMode)
	} else if environment.Active().IsPro() {
		// 生成环境, 将gin的模式, 设置为发布模式
		gin.SetMode(gin.ReleaseMode)
	} else if environment.Active().IsPre() {
		// 预发布环境, 将gin的模式, 设置为发布模式
		gin.SetMode(gin.ReleaseMode)
	} else {
		// 自定义环境, 将gin的模式, 设置为Debug模式
		gin.SetMode(gin.DebugMode)
	}

	// 关闭gin的默认验证器, 一般不要关闭.
	//gin.DisableBindValidation()
	m := &mux.Mux{
		Engine: gin.New(),
	}

	// 执行选项处理函数, 用来控制一些选项开关
	opt := new(mux.Option)
	for _, f := range options {
		f(opt)
	}

	// 没有显式的关闭PProf, 则默认开启PProf
	if !opt.DisablePProf {
		// 加一重双保险, 如果是生产环境的话, 永远不会开启性能分析工具PProf, 避免留下漏洞.
		if !environment.Active().IsPro() {
			pprof.Register(m.Engine) // register pprof to gin
			fmt.Println("* [register pprof]")
		}
	}

	// 开启普罗米修斯
	if !opt.DisablePrometheus {
		m.Engine.GET("/metrics", gin.WrapH(promhttp.Handler())) // register prometheus
		fmt.Printf("* [register prometheus %s/metrics]\n", configs.Default().Register.Url())
	}

	// 没有显式的关闭swagger, 则默认开启swagger
	if !opt.DisableSwagger {
		// 加一重双保险, 如果是生产环境的话, 永远不会开启swagger, 避免留下漏洞.
		// swagger 开启后访问uri /swagger/index.html 即可访问swagger接口
		if !environment.Active().IsPro() {
			m.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // register swagger
			fmt.Printf("* [register swagger %s/swagger/index.html]\n", configs.Default().Register.Url())
		}
	}

	// 开启cors跨域资源共享
	SetCors(opt, m)

	// 第二次recover(第一次在初始化mux.Context时)，防止处理时发生panic，尤其是在OnPanicNotify中。
	m.Engine.Use(func(ctx *gin.Context) {
		defer func() {
			if gin.Mode() != gin.DebugMode {
				if err := recover(); err != nil {
					logger.Error("got panic", zap.String("panic", fmt.Sprintf("%+v", err)), zap.String("stack", string(debug.Stack())))
					// 程序执行过程中发生panic, 抛出 500 状态码
					ctx.JSON(http.StatusInternalServerError, &businessCodex.Response{
						Code: businessCodex.GetServerErrorCode(),
						Msg:  businessCodex.Text(businessCodex.GetServerErrorCode()),
					})
				}
			}

		}()

		ctx.Next()
	})
	// 在中间件中, 初始化mux.Context
	m.Engine.Use(mux.InitContext(*r, *opt))

	// 限速器
	SetLimiter(opt, m)

	m.Engine.NoMethod(mux.WrapHandlers(mux.DisableTrace)...)
	m.Engine.NoRoute(mux.WrapHandlers(mux.DisableTrace)...)
	// 设置系统接口
	SetSystem(m)
	if defaultMux == nil {
		defaultMux = m
	}
	return m, nil
}

func SetCors(opt *mux.Option, m *mux.Mux) {
	// 开启cors跨域资源共享
	if opt.EnableCors {
		m.Engine.Use(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:     []string{"*"},
			AllowCredentials:   true,
			OptionsPassthrough: true,
		}))
	}
}

func SetLimiter(opt *mux.Option, m *mux.Mux) {
	if opt.EnableRate {
		// 开启限速器
		limiter := rate.NewLimiter(rate.Every(time.Second*1), mux.MaxBurstSize)
		m.Engine.Use(func(ctx *gin.Context) {
			context := mux.NewContext(ctx)
			defer mux.ReleaseContext(context)

			if !limiter.Allow() {
				context.AbortWithError(errno.New429Errno(businessCodex.GetTooManyRequestsCode(),
					errno.NewError(businessCodex.Text(businessCodex.GetTooManyRequestsCode()))),
				)
				return
			}

			ctx.Next()
		})
	}
}

type HealthResponse struct {
	Version      string `json:"version"`
	BuildVersion string `json:"buildVersion"`
	BuildDate    string `json:"buildDate"`
	Datetime     string `json:"datetime"`
	Environment  string `json:"environment"`
	Service      string `json:"service"`
	Host         string `json:"host"`
	Status       string `json:"status"`
}

func NewHealthResponse() *HealthResponse {
	resp := &HealthResponse{
		Status: "ok",
	}
	return resp
}

func SetSystem(m *mux.Mux) {
	system := m.Group("/system")
	{
		// 健康检查 @router /system/health response {"Status": "ok"}
		system.GET("/health", func(ctx mux.Context) {
			r := NewHealthResponse()
			config := configs.Default()
			r.BuildVersion = config.BuildInfo.BuildVersion
			r.BuildDate = config.BuildInfo.BuildDate
			r.Version = configs.Default().Register.ProjectVersion
			r.Service = configs.Default().Register.Name
			r.Datetime = timex.JSONTimeNow().String()
			r.Environment = environment.Active().Value()
			r.Host = ctx.Host()
			ctx.Payload(r)
		})

		fmt.Printf("* [register health %s/system/health] \n", configs.Default().Register.Url())
	}
}
