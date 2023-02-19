package ginWarpExampleApi

import (
	"net/http"
	"testing"

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

func TestUserDetail(t *testing.T) {
	ctx := getTestContex()
	httpDiscover.CreateSpecialHttpClient()
	type args struct {
		ctx     mux.Context
		request UserDetailRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"account1", args{ctx: ctx, request: UserDetailRequest{
			Account: "admin1",
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, code, err := UserDetail(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if code != http.StatusOK {
				t.Errorf("UserDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(resp)
		})
	}
}
