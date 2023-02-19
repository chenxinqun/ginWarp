package usersRepo

import (
	"github.com/chenxinqun/ginWarpPkg/cryptox/password"
	"github.com/chenxinqun/ginWarpPkg/timex"
)

func (t *Users) SetPassword(inputPasswd string) {
	t.Password = password.GeneratePassword(inputPasswd)
}

func (t *Users) SetLoginTime() {
	t.LoginTime = timex.JSONTimeNow()
}
