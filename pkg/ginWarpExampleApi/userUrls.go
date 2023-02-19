package ginWarpExampleApi

import "github.com/chenxinqun/ginWarp/pkg/urls"

// 路由组相关常量

const (
	UserGroup = "/user"
)

// user路由组相关URL

var (
	V1UserCreateUrl urls.Url
	V1UserUpdateUrl urls.Url
	V1UserDeleteUrl urls.Url
	V1UserDetailUrl urls.Url
)

func init() {
	V1UserCreateUrl = urls.New().R(RootGroup).S(ServiceName).V(V1).G(UserGroup).P("/create")
	V1UserUpdateUrl = urls.New().R(RootGroup).S(ServiceName).V(V1).G(UserGroup).P("/update")
	V1UserDeleteUrl = urls.New().R(RootGroup).S(ServiceName).V(V1).G(UserGroup).P("/delete")
	V1UserDetailUrl = urls.New().R(RootGroup).S(ServiceName).V(V1).G(UserGroup).P("/detail")
}
