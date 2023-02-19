# ginWarp

一个基于gin框架二次封装的框架

# 技术背景

ginWarp 框架是基于 [gin框架](https://markdown.com.cn) 的二开

主要参考项目是gin的二开项目 [go-gin-api](https://www.yuque.com/xinliangnote/go-gin-api) 

底层基于gin和gorm. 相关问题都可以找到开源方案. 我只整改了上层代码, 改变了使用方式.

相对  [go-gin-api](https://www.yuque.com/xinliangnote/go-gin-api) ,主要不同点在于,这是一个非侵入式框架. 

同时做了许多裁剪, 包含的内容不如 [go-gin-api](https://www.yuque.com/xinliangnote/go-gin-api) 全面.

# ginWarp

### Go Get 使用 ginWarp

``` bash

# 包名 github.com/chenxinqun/ginWarp
go get -u github.com/chenxinqun/ginWarp
```

### ginWarp 程序启动流程与代码开发流程

1. main函数引用 github.com/chenxinqun/ginWarp/coreServer 包, 调用 coreServer.StartHttpServer(nil, router.SetApiRouter), 传入自己写的
   router.SetApiRouter 函数注册路由
2. StartHttpServer 会内置两个命令行参数, 一个 env, 用来指定环境变量加载不同的配置文件, 默认为dev开发环境; 一个 addr, 传入一个带冒号的端口号, 指定监听地址, 默认为 :8888.
3. StartHttpServer 会调用server/init.go 中的一些列初始化函数, 对配置文件以及各种数据库连接进行初始化
4. 初始化完数据库之后, 调用传入的 router.SetApiRouter 函数, 注册路由
5. 开发人员自己在 router.SetApiRouter 中, 将 handlers中的HandlerFunc注册成URL.
6. 开发人员在 ./internal/api 中根据MVC分层自己实现业务逻辑.
7. ginWarp 提供了MVC的代码生成工具, 可以参考 github.com/chenxinqun/zhcorectl 的使用文档
8. 所有的Repo连接, 都有一个Default()函数, 可以对默认Repo连接进行获取, 在service中可以引用相应的repo所在包进行使用. 如: mysql.Default(), redis.Default(),
   mongo.Default()

#### ginWarp/example/api 目录中有一个简单的MVC示例, 包含上传文件, MySQL的修改查询与创建. 可以作为参考. example的main函数在 cmd/example中

# 工程目录结构

``` text

ginWarp
├─cmd                            // 存放一些main文件
│  ├─example                     // ginWarp/example 的启动目录
│  └─pushconfig                  // 启动 ginWarp/example 所需的配置文件
├─configs                        // 配置文件以及处理配置项的代码
├─docs                           // 文档目录
├─example                        // 一些框架使用的简单例子, 对等internal/目录.
│  ├─api                         
│  │  ├─handlers                 
│  │  │  └─userHandler           
│  │  ├─repository               
│  │  │  ├─mongoRepo             
│  │  │  ├─mysqlRepo             
│  │  │  │  └─usersRepo          
│  │  │  └─redisRepo             
│  │  └─services                 
│  │      └─userService          
│  └─router                     
├─internal                       // 业务内部包, 项目内部业务逻辑实现(业务代码)
│  ├─api                         // 存放 mvc 代码
│  │  ├─handlers                 // handlers层
│  │  ├─repository               // model层
│  │  └─services                 // service层
│  ├─pkg                         // 业务相关公共包
│  └─router                      // 路由注册
├─middleware                     // 一些全局的中间件
├─pkg                            // 公共包, 与业务关联较小,或多个业务可用,以及可以被别的项目所引用的包
│  ├─core                        // 关于Gin引擎的一些封装, 以及内置路由的一些定义
│  ├─coreServer                  // http服务启动与停止的一些封装
│  ├─ginWarpExampleApi           // 关于ginWarp/example项目用到的url的一些定义.
│  ├─ginWarpExampleCode          // 示例代码 ginWarp/example 用到的一些URL, request, response以及Client的定义
│  ├─ginWarpExampleRedis         // 示例代码 ginWarp/example 用到的一些业务码
│  ├─session                     // 示例代码 ginWarp/example 用到的 关于redis调用的一些封装
│  └─urls                        // 关于url的公共定义
├─scripts                        // 存放一些脚本文件
├─server                         // 具体业务中才需要的目录, 主要是定义一些在coreServer包中所用到的init与close方法
└─testing                        // 黑盒测试的测试用例存放目录
             
```

# 功能使用教程

### 参数验证器

参数验证器已经在框架中集成了. 只需要在接收参数的结构体中增加响应的Tag就行了. 关键字使用`binding:""`
需要先之的内容, 定义在binding:""之中. 如下图所示:

```go

type CreateRequest struct {
Account string `json:"account" binding:"required,min=4,max=24"`
Pwd     string `json:"pwd" binding:"required,min=8,max=16"`
}

// binding是一个适用于Gin框架的开源验证器. 底层使用 github.com/go-playground/validator 验证器.
```

##### 具体定义字段, 可以参考 github.com/go-playground/validator 的 [文档](https://github.com/go-playground/validator) .

### 使用MySQL

MySQL连接使用全局变量, 在 github.com/chenxinqun/ginWarpPkg/datax/mysqlx 包中. 用以下函数可以获得:
mysqlx.Default()

MySQL连接需要配合代码生成器生成的DAO层使用. 所有的传入DAO层执行SQL的conn, 统一使用 mysql.Default()获取.

如果要获取原生的gorm.DB 请使用 mysql.Default().GetDb().WithContext(ctx) 方法. 示例代码如下:

```go
package userService

import (
	"github.com/chenxinqun/ginWarp/example/api/repository/mysqlRepo/usersRepo"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleApi"
	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleCode"

	"github.com/chenxinqun/ginWarpPkg/convert"
	"github.com/chenxinqun/ginWarpPkg/datax/mysqlx"
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
)

type CreateParams struct {
	Account string
	Passwd  string
}

type CreateResult struct {
	User *usersRepo.Users
}

// Create 创建用户
// @Author Test
// @Summary 创建用户
// @Description 创建用户
// @Handlers handlers/asset/userHandler
func (s *service) Create(ctx mux.Context, params *ginWarpExampleApi.UserCreateRequest) (ret *ginWarpExampleApi.UserCreateResponse, code int, err error) {
	ret = new(ginWarpExampleApi.UserCreateResponse)
	user := usersRepo.NewModel()
	user.Account = params.Account
	user.SetPassword(params.Pwd)
	user.SetLoginTime()
	// 全局MySQL连接
	repo := mysqlx.Default()
	// 已经封装了超时时间的请求专用context
	rctx := ctx.RequestContext()
	_, err = user.Create(rctx, repo)
	if err != nil {
		return nil, ginWarpExampleCode.UserCreateError, err
	}
	err = convert.StructToStruct(*user, ret)
	if err != nil {
		return nil, ginWarpExampleCode.UserCreateError, err
	}
	return
}


```

MySQL事务:
使用 dao.Transaction 传入一个ctx, 一个conn和一个事务函数. 如下所示:

```go
dao := usersRepo.NewQueryBuilder()
// 全局MySQL连接
repo := mysqlx.Default()
// 已经封装了超时时间的请求专用context
rctx := ctx.RequestContext()
dao.Transaction(rctx, repo, func (tx *gorm.DB) error {
// 在事务中执行一些 db 操作（从这里开始，您应该使用 'mysql' 而不是 'db'）
if err := tx.Create(&Animal{Name: "Giraffe"}).Error; err != nil {
// 返回任何错误都会回滚事务
return err
}

if err := tx.Create(&Animal{Name: "Lion"}).Error; err != nil {
return err
}

// 返回 nil 提交事务
return nil
})

```

### kafka的使用

MySQL连接使用全局变量, 在 github.com/chenxinqun/ginWarpPkg/datax/kafkax 包中. 用一下四个函数, 可以分别获得:
消费者: kafka.DefaultConsumer()
消费组: kafka.DefaultConsumerGroup()
异步生产者: kafka.DefaultAsyncProducer()
同步生产者: kafka.DefaultProducer()

使用示例如下:

```go
package main

import (
	"github.com/chenxinqun/ginWarpPkg/datax/kafkax"
	"github.com/chenxinqun/ginWarpPkg/errno"
	"github.com/chenxinqun/ginWarpPkg/sysx/environment"
)

type Topic struct {
	kafkax.Topic
	t string
}

func (t *Topic) String() (topic string) {
	topic = t.t
	return
}

func main() {
	t := &Topic{t: "test1"}
	kafkax.DefaultConsumerGroup().Consume([]kafkax.Topic{t}, func(msg *kafkax.ConsumerMessage) (error, bool) {
		if msg == nil {
			return errno.NewError("msg为空"), true
		}
		return nil, true
	})

}

```

### Redis使用

MongoDB连接使用全局变量, 在 github.com/chenxinqun/ginWarpPkg/datax/redisx 包中. 用以下函数可以获得:
redisx.Default()
任何操作都应该使用redis的interface封装. 使用示例如下:

```go
// Get Set string
const (
DetailsKey = "details"
)

func UserCacheKey(userID int64) string {
return fmt.Sprintf("%s:%d", scrm_system_redis.CachePrefix, userID)
}

func SetUserDetails(ctx mux.Context, userID int64, val scrm_user_repo.ScrmUser) error {
key := UserCacheKey(userID)
data, err := json.Marshal(val)
if err != nil {
return err
}
return redisx.Default().HSet(ctx.RequestContext(), key, map[string]interface{}{DetailsKey: string(data)}).Err()
}

func GetUserDetails(ctx mux.Context, userID int64, result *scrm_system_api.UsersDetailsResponse) error {
key := UserCacheKey(userID)
data, err := redisx.Default().HGet(ctx.RequestContext(), key, DetailsKey).Result()
if err != nil {
return err
}
err = json.Unmarshal([]byte(data), result)
return err
}

```

```go
// 发布
func (s *service) PostProve(ctx mux.Context, rawData []byte) (ret map[string]string, err error) {
resp := make(map[string]string)
resp["RawData"] = string(rawData)
respJson, err := json.Marshal(resp)
if err != nil {
return nil, err
}
err = redisx.Default().Publish(string(respJson))
if err != nil {
return nil, err
}
ret = resp
return
}
```

```go

// 订阅
// SubscribeCallBackEvent 企业微信事件回调
func SubscribeCallBackEvent() (err error) {
err = redis.Default().Subscribe(func (msg *redis.Message) (e error) {
ret := make(map[string]string)
// 原始数据再msg.Payload里面
e = json.Unmarshal([]byte(msg.Payload), &ret)
if e != nil {
loggerx.Default().Error("json序列化时报错", zap.Error(e), zap.String("redis数据", msg.Payload))
return e
}
...
下面写业务
})
return
}

```