package ginWarpExampleApi

type UserCreateRequest struct {
	Account string `json:"account" binding:"required,min=4,max=24"` // 用户账号
	Pwd     string `json:"pwd" binding:"required,min=8,max=16"`     // 用户密码
}

type UserDeleteRequest struct {
	ID      int64  `json:"id,string" form:"id"`                                    // 用户ID int64的字段, 统一在json tag中增加 string标识, 避免前端接收数据失真的问题.
	Account string `json:"account" form:"account" binding:"required,min=4,max=24"` // 用户账号
}

type UserDetailRequest struct {
	ID      int64  `json:"id,string" form:"id"`                                    // 用户ID int64的字段, 统一在json tag中增加 string标识, 避免前端接收数据失真的问题.
	Account string `json:"account" form:"account" binding:"required,min=4,max=24"` // 用户账号
}

type UserUpdatePasswordRequest struct {
	Account string `form:"account" json:"account"` // 用户账号
	Pwd     string `form:"pwd" json:"pwd"`         // 用户密码
	NewPwd  string `form:"new_pwd" json:"new_pwd"` // 用户新密码
}
