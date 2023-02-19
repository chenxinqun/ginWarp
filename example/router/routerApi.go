package router

import (
	"github.com/chenxinqun/ginWarp/example/api/handlers/userHandler"
	"github.com/chenxinqun/ginWarp/pkg/core"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleApi"
)

func SetApiRouter(r *core.RouterResource) {
	api := r.Mux.Group(ginWarpExampleApi.RootGroup)
	service := api.Group(core.GetServiceNameUrl())
	v1 := service.Group(ginWarpExampleApi.V1)
	{
		userHd := userHandler.New()
		userGroup := v1.Group(ginWarpExampleApi.UserGroup)
		// 创建
		userGroup.POST(ginWarpExampleApi.V1UserCreateUrl.Path(), userHd.Create())
		// 更新
		userGroup.PUT(ginWarpExampleApi.V1UserUpdateUrl.Path(), userHd.UpdatePassword())
		// 删除
		userGroup.DELETE(ginWarpExampleApi.V1UserDeleteUrl.Path(), userHd.Delete())
		// 查询
		userGroup.GET(ginWarpExampleApi.V1UserDetailUrl.Path(), userHd.Detail())
	}
}
