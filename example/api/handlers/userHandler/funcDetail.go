package userHandler

import (
	"fmt"

	"github.com/chenxinqun/ginWarp/example/api/services/userService"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleApi"

	"github.com/chenxinqun/ginWarpPkg/errno"
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
	logger "github.com/chenxinqun/ginWarpPkg/loggerx"
	"go.uber.org/zap"
)

type swaggerDetailSuccess struct {
	Code int                                  `json:"code"` // 业务码
	Msg  string                               `json:"msg"`  // 描述信息
	Data ginWarpExampleApi.UserDetailResponse `json:"data"` // 返回值
}

type swaggerDetailFailure struct {
	Code int    `json:"code"` // 业务码
	Msg  string `json:"msg"`  // 描述信息
}

// Detail 用户详情
// @Author Test
// @Summary 用户详情
// @Description 用户详情
// @Tags user
// @Accept json
// @Produce json
// @Param Request query ginWarpExampleApi.UserDetailRequest true "请求信息"
// @Success 200 {object} swaggerDetailSuccess
// @Response 202 {object} swaggerDetailFailure
// @Router /v1/user/detail [get]
func (h *handler) Detail() mux.HandlerFunc {
	return func(c mux.Context) {
		req := new(ginWarpExampleApi.UserDetailRequest)
		res := new(ginWarpExampleApi.UserDetailResponse)
		fmt.Println("UserID", c.UserID())
		fmt.Println("UserName", c.UserName())
		fmt.Println("TenantID", c.TenantID())
		fmt.Println("IsAdmin", c.IsAdmin())
		fmt.Println("RoleType", c.RoleType())
		// 绑定参数的方法未必是生成的这个, 请根据实际情况更换.
		if err := c.ShouldBindQuery(req); err != nil {
			// 全局日志
			logger.Default().Error("参数序列化错误", zap.Error(err))
			c.AbortWithError(err)

			return
		}
		service := userService.New()
		ret, code, err := service.Detail(c, req)
		if err != nil {
			// 调用全局日志记录异常
			logger.Default().Error("查询用户详情错误", zap.Error(err))
			c.AbortWithError(errno.NewBusinessErrno(code, err))
			return
		}
		res = ret
		c.Payload(res)
	}
}
