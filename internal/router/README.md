## internal/router

内部路由包。

配置路由的地方

命名规范：

- 文件名以 router_api.go命名
- router_api.go 内必须有 func SetApiRouter(r *Resource) 函数

```golang
package router

import fc_router "github.com/chenxinqun/ginWarp/example/router"

func SetApiRouter(r *fc_router.Resource) {
    
}

```
