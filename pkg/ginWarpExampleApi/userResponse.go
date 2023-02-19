package ginWarpExampleApi

import "time"

type UserCreateResponse struct {
	ID      int64  `json:"id,string"` // 用户ID int64的字段, 统一在json tag中增加 string标识, 避免前端接收数据失真的问题.
	Account string `json:"account"`   // 用户账号
}

type UserDeleteResponse struct {
	DeleteUser int64 `json:"delete_user,string"` // 如果成功则返回被删除的用户ID, 如果失败则为0
} // 不需要返回值, 通过http状态码就可以判断了

type UserDetailResponse struct {
	ID        int64     `json:"id,string"` // 用户ID int64的字段, 统一在json tag中增加 string标识, 避免前端接收数据失真的问题.
	Account   string    // 用户账号
	LoginTime time.Time `json:"login_time" gorm:"time"` //用户登陆时间
	LoginIp   string    `json:"login_ip"`               // 用户登陆IP
	CreatedAt time.Time `json:"created_at"`             // 用户创建时间
}
type UserUpdatePasswordResponse struct {
	SetPassword bool `json:"set_password"` // 设置密码是否成功
}
