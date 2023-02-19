package userService

import (
	"github.com/chenxinqun/ginWarp/example/api/repository/mysqlRepo/usersRepo"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleApi"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleCode"

	"github.com/chenxinqun/ginWarpPkg/datax/mysqlx"
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type UpdatePasswordParams struct {
	User      usersRepo.Users
	NewPasswd string
}

type UpdatePasswordResult struct {
	User *usersRepo.Users
}

// UpdatePassword 修改用户密码
// @Author Test
// @Summary 创建用户
// @Description 编辑用户
// @Handlers handlers/asset/userHandler
func (s *service) UpdatePassword(ctx mux.Context, params *ginWarpExampleApi.UserUpdatePasswordRequest) (ret *ginWarpExampleApi.UserUpdatePasswordResponse, code int, err error) {
	ret = new(ginWarpExampleApi.UserUpdatePasswordResponse)
	dao := usersRepo.NewQueryBuilder(ctx.TenantID())
	user := usersRepo.NewModel()
	user.SetPassword(params.Pwd)
	dao = dao.WhereAccount(mysqlx.EqualPredicate, params.Account)
	dao = dao.WherePassword(mysqlx.EqualPredicate, user.Password)
	// 全局MySQL连接
	repo := mysqlx.Default()
	// 已经封装了超时时间的请求专用context
	rctx := ctx.RequestContext()
	exist, _ := dao.Exist(rctx, repo)
	if !exist {
		return nil, ginWarpExampleCode.UserNotFoundError, errors.New("账号不存在或密码错误")
	}
	user.SetPassword(params.NewPwd)
	rctx = ctx.RequestContext()
	err = dao.Updates(rctx, repo, gin.H{"password": user.Password})
	if err != nil {
		return nil, ginWarpExampleCode.UserUpdateError, err
	}
	ret.SetPassword = true
	return
}
