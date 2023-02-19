package userService

import (
	"github.com/chenxinqun/ginWarp/example/api/repository/mysqlRepo/usersRepo"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleApi"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleCode"

	"github.com/chenxinqun/ginWarpPkg/convert"
	"github.com/chenxinqun/ginWarpPkg/datax/mysqlx"
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
)

type CreateParams struct {
	Account string
	Passwd  string
}

type CreateResult struct {
	User *usersRepo.Users
}

// Create 创建用户
// @Author Test
// @Summary 创建用户
// @Description 创建用户
// @Handlers handlers/asset/userHandler
func (s *service) Create(ctx mux.Context, params *ginWarpExampleApi.UserCreateRequest) (ret *ginWarpExampleApi.UserCreateResponse, code int, err error) {
	ret = new(ginWarpExampleApi.UserCreateResponse)
	user := usersRepo.NewModel()
	user.Account = params.Account
	user.SetPassword(params.Pwd)
	user.SetLoginTime()
	// 全局MySQL连接
	repo := mysqlx.Default()
	// 已经封装了超时时间的请求专用context
	rctx := ctx.RequestContext()
	_, err = user.Create(rctx, repo)
	if err != nil {
		return nil, ginWarpExampleCode.UserCreateError, err
	}
	err = convert.StructToStruct(*user, ret)
	if err != nil {
		return nil, ginWarpExampleCode.UserCreateError, err
	}
	return
}
