package middleware

import (
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
)

func (m *Middleware) DisableLog() mux.HandlerFunc {
	return func(c mux.Context) {
		mux.DisableTrace(c)
	}
}
