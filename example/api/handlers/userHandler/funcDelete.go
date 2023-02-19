package userHandler

import (
	"github.com/chenxinqun/ginWarp/example/api/services/userService"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleApi"

	"github.com/chenxinqun/ginWarpPkg/errno"
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
	logger "github.com/chenxinqun/ginWarpPkg/loggerx"
	"go.uber.org/zap"
)

type swaggerDeleteSuccess struct {
	Code int                                  `json:"code"` // 业务码
	Msg  string                               `json:"msg"`  // 描述信息
	Data ginWarpExampleApi.UserDeleteResponse `json:"data"` // 返回值
}

type swaggerDeleteFailure struct {
	Code int    `json:"code"` // 业务码
	Msg  string `json:"msg"`  // 描述信息
}

// Delete 删除用户
// @Author Test
// @Summary 删除用户
// @Description 删除用户
// @Tags user
// @Accept json
// @Produce json
// @Param Request query ginWarpExampleApi.UserDeleteRequest true "请求信息"
// @Success 200 {object} swaggerDeleteSuccess
// @Response 202 {object} swaggerDeleteFailure
// @Router /v1/user/delete [delete]
func (h *handler) Delete() mux.HandlerFunc {
	return func(c mux.Context) {
		req := new(ginWarpExampleApi.UserDeleteRequest)
		res := new(ginWarpExampleApi.UserDeleteResponse)
		// 绑定参数的方法未必是生成的这个, 请根据实际情况更换.
		if err := c.ShouldBindQuery(req); err != nil {
			c.AbortWithError(err)

			return
		}

		// service相关的调用写在这里.
		service := userService.New()
		ret, code, err := service.Delete(c, req)
		if err != nil {
			// 调用全局日志记录异常
			logger.Default().Error("删除用户错误", zap.Error(err))
			c.AbortWithError(errno.NewBusinessErrno(code, err))

			return
		}
		res = ret
		c.Payload(res)
	}
}
