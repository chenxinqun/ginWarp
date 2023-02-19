package ginWarpExampleCode

import "github.com/chenxinqun/ginWarpPkg/businessCodex"

// 实际业务使用六位数的业务码. 前两位表示微服务. 中间两位表示业务模块, 末两位表示具体错误类型
const (
	// SucceedCode 成功业务码, 成功业务码不使用0, 是为了避免在一些不需要返回值的接口中. 用以区分是否真的成功了.
	SucceedCode = 1000000000
)

func init() {
	businessCodex.SetSucceedCode(SucceedCode)
	businessCodex.Init(false)
}

// 这些是zh_example的示例代码用到的业务码, 为了不跟实际业务发生冲突, 因此使用五位数的业务码
const (
	UserCreateError       = 20201
	UserDeleteError       = 20203
	UserUpdateError       = 20204
	UserDetailError       = 20213
	UserNotFoundError     = 20214
	UserUploadAvatarError = 20215
	UserSaveAvatarError   = 20216
)
