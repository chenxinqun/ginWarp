package userService

import (
	"github.com/chenxinqun/ginWarp/example/api/repository/mysqlRepo/usersRepo"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleApi"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleCode"

	"github.com/chenxinqun/ginWarpPkg/datax/mysqlx"
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
)

type DeleteParams struct{}

type DeleteResult struct{}

// Delete 删除用户
// @Author Test
// @Summary 删除用户
// @Description 删除用户
// @Handlers handlers/asset/userHandler
func (s *service) Delete(ctx mux.Context, params *ginWarpExampleApi.UserDeleteRequest) (ret *ginWarpExampleApi.UserDeleteResponse, code int, err error) {
	ret = new(ginWarpExampleApi.UserDeleteResponse)

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
	err = dao.Delete(rctx, repo)
	if err != nil {
		return nil, ginWarpExampleCode.UserDeleteError, err
	}
	ret.DeleteUser = params.ID
	return
}
