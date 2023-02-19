package ginWarpExampleApi

import (
	"github.com/chenxinqun/ginWarpPkg/httpx/httpDiscover"
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
)

func UserDetail(ctx mux.Context, request UserDetailRequest) (UserDetailResponse, int, error) {
	ret := new(UserDetailResponse)
	client := httpDiscover.New()
	code, err := client.GetJson(ctx, ServiceName, V1UserDetailUrl.String(), request, ret)
	return *ret, code, err
}
