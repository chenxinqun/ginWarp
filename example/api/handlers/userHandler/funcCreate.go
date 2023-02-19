package userHandler

import (
	"github.com/chenxinqun/ginWarp/example/api/services/userService"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleApi"

	"github.com/chenxinqun/ginWarpPkg/errno"
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
	logger "github.com/chenxinqun/ginWarpPkg/loggerx"
	"go.uber.org/zap"
)

type swaggerCreateSuccess struct {
	Code int                                  `json:"code"` // 业务码
	Msg  string                               `json:"msg"`  // 描述信息
	Data ginWarpExampleApi.UserCreateResponse `json:"data"` // 返回值
}
type swaggerCreateFailure struct {
	Code int    `json:"code"` // 业务码
	Msg  string `json:"msg"`  // 描述信息
}

// Create 创建用户
// @Author Test
// @Summary 创建用户
// @Description 创建用户
// @Tags user
// @Accept json
// @Produce json
// @Param Request body ginWarpExampleApi.UserCreateRequest true "请求信息"
// @Success 200 {object} swaggerCreateSuccess
// @Response 202 {object} swaggerCreateFailure
// @Router /v1/user/create [post]
func (h *handler) Create() mux.HandlerFunc {
	return func(c mux.Context) {
		req := new(ginWarpExampleApi.UserCreateRequest)
		res := new(ginWarpExampleApi.UserCreateResponse)
		// 绑定参数的方法未必是生成的这个, 请根据实际情况更换.
		if err := c.ShouldBindJSON(req); err != nil {
			// 全局日志
			logger.Default().Error("参数序列化错误", zap.Error(err))
			c.AbortWithError(err)

			return
		}
		// service相关的调用写在这里.
		service := userService.New()
		ret, code, err := service.Create(c, req)
		if err != nil {
			// 调用全局日志记录异常
			logger.Default().Error("创建用户错误", zap.Error(err))
			c.AbortWithError(errno.NewBusinessErrno(code, err))

			return
		}
		res = ret
		c.Payload(res)
	}
}
