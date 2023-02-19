package main

import (
	_ "github.com/chenxinqun/ginWarp/docs"
	"github.com/chenxinqun/ginWarp/example/router"
	"github.com/chenxinqun/ginWarp/pkg/coreServer"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleApi"
	"github.com/chenxinqun/ginWarp/server"
	"log"
)

// @title ginWarp-example 接口文档
// @version v1
// @description
// @BasePath /api/ginWarp-example/

var (
	BuildDate    string // 编译时,由CI/CD工具复写此变量, 并赋值给coreServer.Resource, 传入coreServer的启动函数中
	BuildVersion string // 编译时,由CI/CD工具复写此变量, 并赋值给coreServer.Resource, 传入coreServer的启动函数中
)

func main() {
	// 解析命令行
	args := coreServer.FlagInit()
	// 定义微服务名称和端口, 这里支持显式定义, 一般用于开发模式, 便于理解.
	// 部署时, 可以通过配置文件, 或者启动命令行, 重新赋值覆盖. 使程序部署形式更加灵活, 便于做统一的部署管理.
	r := &coreServer.Resource{
		BuildVersion: BuildVersion, BuildDate: BuildDate,
		ServiceName: ginWarpExampleApi.ServiceName, Addr: ":8888"}
	// 注册自定义初始化函数
	coreServer.RegisterInitFunc(server.InitExample)
	// 注册自定义关闭函数
	coreServer.RegisterCloseFunc(server.CloseExample)
	// 非http模式示例
	if args.NoHttpMod() {
		coreServer.StartHttpServer(r, func() error {
			log.Fatalln("这是一个示例脚本")
			return nil
		})
	} else {
		// 注册路由并启动服务
		coreServer.StartHttpServer(r, nil, router.SetApiRouter)
		// 监听关闭信号并关闭服务
		coreServer.CloseCoreServer()
	}
}
