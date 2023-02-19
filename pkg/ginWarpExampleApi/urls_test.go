package ginWarpExampleApi

import (
	"testing"
)

func TestUrlDefine(t *testing.T) {
	t.Log(V1UserCreateUrl.String())
	t.Log(V1UserUpdateUrl.String())
	t.Log(V1UserDeleteUrl.String())
	t.Log(V1UserDetailUrl.String())
}
