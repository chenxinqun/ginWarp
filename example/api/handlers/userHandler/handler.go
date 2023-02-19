package userHandler

import (
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
)

var _ Handler = (*handler)(nil)

type Handler interface {

	// Create 创建用户
	// @Author Test
	// @Uri /v1/user/create [post]
	Create() mux.HandlerFunc

	// UpdatePassword 编辑用户
	// @Author Test
	// @Uri /v1/user/update [put]
	UpdatePassword() mux.HandlerFunc

	// Delete 删除用户
	// @Author Test
	// @Uri /v1/user/delete delete]
	Delete() mux.HandlerFunc

	// Detail 用户详情
	// @Author Test
	// @Uri /v1/user/detail [get]
	Detail() mux.HandlerFunc
}

type handler struct{}

func New() Handler {
	return &handler{}
}
