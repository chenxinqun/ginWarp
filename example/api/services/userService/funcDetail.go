package userService

import (
	"github.com/chenxinqun/ginWarp/example/api/repository/mysqlRepo/usersRepo"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleApi"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleCode"

	"github.com/chenxinqun/ginWarpPkg/convert"
	"github.com/chenxinqun/ginWarpPkg/datax/mysqlx"
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
)

type DetailParams struct {
	Account string
	Passwd  string
}

type DetailResult struct {
	User *usersRepo.Users
}

// Detail 用户详
// @Author Test
// @Summary 用户详情
// @Description 用户详情
// @Handlers handlers/asset/userHandler
func (s *service) Detail(ctx mux.Context, params *ginWarpExampleApi.UserDetailRequest) (ret *ginWarpExampleApi.UserDetailResponse, code int, err error) {
	ret = new(ginWarpExampleApi.UserDetailResponse)

	dao := usersRepo.NewQueryBuilder(ctx.TenantID())
	if params.ID > 0 {
		dao = dao.WhereID(mysqlx.EqualPredicate, params.ID)
	} else {
		dao = dao.WhereAccount(mysqlx.EqualPredicate, params.Account)
	}
	// 全局MySQL连接
	repo := mysqlx.Default()
	// 已经封装了超时时间的请求专用context
	rctx := ctx.RequestContext()
	user, err := dao.First(rctx, repo)
	if err != nil {
		return nil, ginWarpExampleCode.UserDetailError, err
	}
	err = convert.StructToStruct(*user, ret)
	if err != nil {
		return nil, ginWarpExampleCode.UserDetailError, err
	}
	return
}
