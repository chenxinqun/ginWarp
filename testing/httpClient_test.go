package testing

import (
	"testing"

	"github.com/chenxinqun/ginWarp/pkg/ginWarpExampleApi"

	"github.com/chenxinqun/ginWarpPkg/httpx/httpDiscover"
	"github.com/chenxinqun/ginWarpPkg/httpx/mux"
)

func getTestContex() mux.Context {
	ctx := mux.CreateSpecialContext(mux.SpecialContextResource{
		UserID:   1234567890123456789,
		UserName: "aaaaaa",
		TenantID: 12345678,
		IsAdmin:  true,
		RoleType: 1,
	})
	return ctx
}

func TestServiceClient_GetJson(t *testing.T) {
	ctx := getTestContex()
	client := httpDiscover.CreateSpecialHttpClient()
	t.Log(ctx.UserID())
	t.Log(ctx.UserName())
	t.Log(ctx.TenantID())
	t.Log(ctx.IsAdmin())
	t.Log(ctx.RoleType())
	req := ginWarpExampleApi.UserDetailRequest{
		Account: "admin1",
	}
	body := new(ginWarpExampleApi.UserDetailResponse)
	httpCode, err := client.GetJson(ctx, ginWarpExampleApi.ServiceName, ginWarpExampleApi.V1UserDetailUrl.String(), req, body)
	if body.ID == 0 || httpCode == 0 || err != nil {
		t.Fatal(body, httpCode, err)
	}
	t.Log(body, httpCode, err)
}

func TestServiceClient_Delete(t *testing.T) {
	ctx := getTestContex()
	client := httpDiscover.CreateSpecialHttpClient()
	t.Log(ctx.UserID())
	t.Log(ctx.UserName())
	t.Log(ctx.TenantID())
	t.Log(ctx.IsAdmin())
	t.Log(ctx.RoleType())
	req := new(ginWarpExampleApi.UserDeleteRequest)
	body := new(ginWarpExampleApi.UserDeleteResponse)
	httpCode, err := client.DeleteJson(ctx, ginWarpExampleApi.ServiceName, ginWarpExampleApi.V1UserDeleteUrl.String(), req, body)
	if body == nil || httpCode == 0 || err != nil {
		t.Fatal(body, httpCode, err)
	}
	t.Log(body, httpCode, err)

}

func TestServiceClient_PostJson(t *testing.T) {
	ctx := getTestContex()
	client := httpDiscover.CreateSpecialHttpClient()
	t.Log(ctx.UserID())
	t.Log(ctx.UserName())
	t.Log(ctx.TenantID())
	t.Log(ctx.IsAdmin())
	t.Log(ctx.RoleType())
	req := ginWarpExampleApi.UserCreateRequest{
		Account: "admin2",
		Pwd:     "admin1234",
	}
	body := new(ginWarpExampleApi.UserCreateResponse)
	httpCode, err := client.PostJson(ctx, ginWarpExampleApi.ServiceName, ginWarpExampleApi.V1UserCreateUrl.String(), req, body)
	if body.ID == 0 || httpCode == 0 || err != nil {
		t.Fatal(body, httpCode, err)
	}
	t.Log(body, httpCode, err)

}

func TestServiceClient_PutJson(t *testing.T) {
	ctx := getTestContex()
	client := httpDiscover.CreateSpecialHttpClient()
	t.Log(ctx.UserID())
	t.Log(ctx.UserName())
	t.Log(ctx.TenantID())
	t.Log(ctx.IsAdmin())
	t.Log(ctx.RoleType())
	req := ginWarpExampleApi.UserUpdatePasswordRequest{}
	body := new(ginWarpExampleApi.UserUpdatePasswordResponse)
	httpCode, err := client.PutJson(ctx, ginWarpExampleApi.ServiceName, ginWarpExampleApi.V1UserUpdateUrl.String(), req, body)
	if body.SetPassword == false || httpCode == 0 || err != nil {
		t.Fatal(body, httpCode, err)
	}
	t.Log(body, httpCode, err)
}
