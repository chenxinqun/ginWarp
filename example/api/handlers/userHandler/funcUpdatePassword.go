package userHandler

import (
	"github.com/chenxinqun/ginWarp/example/api/services/userService"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleApi"

	"github.com/chenxinqun/ginWarpPkg/errno"
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
	logger "github.com/chenxinqun/ginWarpPkg/loggerx"
	"go.uber.org/zap"
)

type swaggerUpdatePasswordSuccess struct {
	Code int                                          `json:"code"` // 业务码
	Msg  string                                       `json:"msg"`  // 描述信息
	Data ginWarpExampleApi.UserUpdatePasswordResponse `json:"data"` // 返回值
}

type swaggerUpdatePasswordFailure struct {
	Code int    `json:"code"` // 业务码
	Msg  string `json:"msg"`  // 描述信息
}

// UpdatePassword 修改用户密码
// @Author Test
// @Summary 修改用户密码
// @Description 修改用户密码
// @Tags user
// @Accept json
// @Produce json
// @Param Request body ginWarpExampleApi.UserUpdatePasswordRequest true "请求信息"
// @Success 200 {object} swaggerUpdatePasswordSuccess
// @Response 202 {object} swaggerUpdatePasswordFailure
// @Router /v1/user/updatePassword [put]
func (h *handler) UpdatePassword() mux.HandlerFunc {
	return func(c mux.Context) {
		req := new(ginWarpExampleApi.UserUpdatePasswordRequest)
		res := new(ginWarpExampleApi.UserUpdatePasswordResponse)
		// 绑定参数的方法未必是生成的这个, 请根据实际情况更换.
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(err)

			return
		}
		// service相关的调用写在这里.
		service := userService.New()
		ret, code, err := service.UpdatePassword(c, req)
		if err != nil {
			// 调用全局日志记录异常
			logger.Default().Error("修改用户密码错误", zap.Error(err))
			c.AbortWithError(errno.NewBusinessErrno(code, err))
			return
		}
		res = ret

		c.Payload(res)
	}
}
