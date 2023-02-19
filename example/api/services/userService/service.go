package userService

import (
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleApi"

	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
)

var _ Service = (*service)(nil)

type Service interface {
	// Create 创建用户
	// @Author Test
	// @Handlers handlers/asset/userHandler
	Create(ctx mux.Context, params *ginWarpExampleApi.UserCreateRequest) (ret *ginWarpExampleApi.UserCreateResponse, code int, err error)

	// UpdatePassword 修改用户密码
	// @Author Test
	// @Handlers handlers/asset/userHandler
	UpdatePassword(ctx mux.Context, params *ginWarpExampleApi.UserUpdatePasswordRequest) (ret *ginWarpExampleApi.UserUpdatePasswordResponse, code int, err error)

	// Delete 删除用户
	// @Author Test
	// @Handlers handlers/asset/userHandler
	Delete(ctx mux.Context, params *ginWarpExampleApi.UserDeleteRequest) (ret *ginWarpExampleApi.UserDeleteResponse, code int, err error)

	// Detail 用户详情
	// @Author Test
	// @Handlers handlers/asset/userHandler
	Detail(ctx mux.Context, params *ginWarpExampleApi.UserDetailRequest) (ret *ginWarpExampleApi.UserDetailResponse, code int, err error)
}

type service struct{}

func New() Service {
	return &service{}
}
