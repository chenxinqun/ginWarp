package usersRepo

import (
	"github.com/chenxinqun/ginWarpPkg/timex"
)

// Users users表
type Users struct {
	ID        int64          `json:"id" gorm:"column:id"`                      // 用户ID
	TenantID  int64          `json:"tenant_id" gorm:"column:tenant_id"`        // 企业ID
	LoginTime timex.JSONTime `json:"login_time" gorm:"column:login_time;time"` // 登陆时间
	CreatedAt timex.JSONTime `json:"created_at" gorm:"column:created_at;time"` // 创建时间
	UpdatedAt timex.JSONTime `json:"updated_at" gorm:"column:updated_at;time"` // 更新时间
	DeletedAt timex.JSONTime `json:"deleted_at" gorm:"column:deleted_at;time"` // 删除时间
	Account   string         `json:"account" gorm:"column:account"`            // 账号
	Password  string         `json:"password" gorm:"column:password"`          // 密码
	LoginIP   string         `json:"login_ip" gorm:"column:login_ip"`          // 登陆IP
}
