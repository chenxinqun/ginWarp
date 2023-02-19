package ginWarpExampleCode

import "github.com/chenxinqun/ginWarpPkg/businessCodex"

func init() {
	businessCodex.SetEnUSText(enUSText)
}

var enUSText = map[int]string{

	UserCreateError:       "Failed to create user",
	UserDeleteError:       "Failed to delete user",
	UserUpdateError:       "Failed to update user",
	UserDetailError:       "Failed to get personal information",
	UserNotFoundError:     "user not found",
	UserUploadAvatarError: "upload avatar error",
	UserSaveAvatarError:   "save avatar error",
}
