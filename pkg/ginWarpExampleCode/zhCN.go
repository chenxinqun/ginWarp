package ginWarpExampleCode

import "github.com/chenxinqun/ginWarpPkg/businessCodex"

func init() {
	businessCodex.SetZhCNText(zhCNText)
}

var zhCNText = map[int]string{

	UserCreateError:       "创建用户失败",
	UserDeleteError:       "删除用户失败",
	UserUpdateError:       "更新用户失败",
	UserDetailError:       "获取个人信息失败",
	UserNotFoundError:     "用户不存在",
	UserUploadAvatarError: "上传用户头像错误",
	UserSaveAvatarError:   "保存用户头像错误",
}
