package usersRepo

import (
	"context"
	"database/sql"
	"github.com/chenxinqun/ginWarpPkg/datax/mysqlx"
	"github.com/chenxinqun/ginWarpPkg/datax/pagex"
	"github.com/chenxinqun/ginWarpPkg/errno"
	"github.com/chenxinqun/ginWarpPkg/idGen"
	"gorm.io/gorm"
	"time"
)

// 必传字段
const (
	tenantIDField = "tenant_id"
)

// 字段常量
const (
	IDField = "id"

	TenantIDField = "tenant_id"

	LoginTimeField = "login_time"

	CreatedAtField = "created_at"

	UpdatedAtField = "updated_at"

	DeletedAtField = "deleted_at"

	AccountField = "account"

	PasswordField = "password"

	LoginIPField = "login_ip"
)

func NewModel() *Users {
	return new(Users)
}

func NewQueryBuilder(tenantID int64) *usersRepoQueryBuilder {
	if tenantID <= 0 {
		panic(errno.Errorf("tenantID 不能小于 0"))
	}
	builder := &usersRepoQueryBuilder{QueryBuilder: *mysqlx.NewQueryBuilder(&Users{}, tenantIDField)}
	builder.Where(mysqlx.EqualPredicate, tenantIDField, tenantID)
	return builder
}

func (t *Users) Create(ctx context.Context, repo mysqlx.Repo) (id int64, err error) {
	if idGen.Default() != nil {
		// 使用ID生成器生成ID
		id, err = idGen.Default().NextID()
		if err != nil {
			return 0, err
		}
		t.ID = id
	}
	db := repo.GetDb().WithContext(ctx)
	rowsAffected, err := NewQueryBuilder(t.TenantID).QueryBuilder.Create(db, t)
	if err != nil {
		return 0, err
	}
	if rowsAffected <= 0 {
		return 0, errno.Errorf("受影响行数为%v", rowsAffected)
	}
	return t.ID, nil
}

type usersRepoQueryBuilder struct {
	mysqlx.QueryBuilder
}

func (qb *usersRepoQueryBuilder) Transaction(ctx context.Context, repo mysqlx.Repo, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	db := repo.GetDb().WithContext(ctx)
	return qb.QueryBuilder.Transaction(db, fc, opts...)
}

func (qb *usersRepoQueryBuilder) Create(ctx context.Context, repo mysqlx.Repo, model *Users) (err error) {
	db := repo.GetDb().WithContext(ctx)
	rowsAffected, err := qb.QueryBuilder.Create(db, model)
	if err != nil {
		return err
	}
	if rowsAffected <= 0 {
		return errno.Errorf("受影响行数为%v", rowsAffected)
	}
	return nil
}

func (qb *usersRepoQueryBuilder) Updates(ctx context.Context, repo mysqlx.Repo, m map[string]interface{}) (err error) {
	db := repo.GetDb().WithContext(ctx)
	rowsAffected, err := qb.QueryBuilder.Updates(db, m)
	if err != nil {
		return err
	}
	if rowsAffected <= 0 {
		return errno.Errorf("受影响行数为%v", rowsAffected)
	}
	return nil
}

func (qb *usersRepoQueryBuilder) Delete(ctx context.Context, repo mysqlx.Repo) (err error) {
	db := repo.GetDb().WithContext(ctx)
	rowsAffected, err := qb.QueryBuilder.Delete(db)
	if err != nil {
		return err
	}
	if rowsAffected <= 0 {
		return errno.Errorf("受影响行数为%v", rowsAffected)
	}
	return nil
}

func (qb *usersRepoQueryBuilder) Count(ctx context.Context, repo mysqlx.Repo) (int64, error) {
	db := repo.GetDb().WithContext(ctx)
	return qb.QueryBuilder.Count(db)
}

func (qb *usersRepoQueryBuilder) Exist(ctx context.Context, repo mysqlx.Repo) (bool, error) {
	db := repo.GetDb().WithContext(ctx)
	return qb.QueryBuilder.Exist(db)
}

func (qb *usersRepoQueryBuilder) Page(ctx context.Context, repo mysqlx.Repo, parma *pagex.Params) ([]*Users, *pagex.Result, error) {
	db := repo.GetDb().WithContext(ctx)
	UsersList := make([]*Users, 0)
	ret, err := qb.QueryBuilder.Page(db, &UsersList, parma)
	if err != nil {
		return nil, nil, err
	}
	return UsersList, ret, err
}

func (qb *usersRepoQueryBuilder) First(ctx context.Context, repo mysqlx.Repo) (*Users, error) {
	db := repo.GetDb().WithContext(ctx)
	ret := new(Users)
	err := qb.QueryBuilder.First(db, ret)
	if err != nil {
		return nil, err
	}
	return ret, err
}

func (qb *usersRepoQueryBuilder) Get(ctx context.Context, repo mysqlx.Repo) (*Users, error) {
	db := repo.GetDb().WithContext(ctx)
	ret := new(Users)
	err := qb.QueryBuilder.Get(db, ret)
	if err != nil {
		return nil, err
	}
	return ret, err
}

func (qb *usersRepoQueryBuilder) List(ctx context.Context, repo mysqlx.Repo) ([]*Users, error) {
	db := repo.GetDb().WithContext(ctx)
	ret := make([]*Users, 0)
	err := qb.QueryBuilder.List(db, ret)
	if err != nil {
		return nil, err
	}
	return ret, err
}

func (qb *usersRepoQueryBuilder) WhereID(p mysqlx.Predicate, value int64) *usersRepoQueryBuilder {
	qb.QueryBuilder.Where(p, IDField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereIDIn(value []int64) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereIn(IDField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereIDNotIn(value []int64) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereNotIn(IDField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByIDASC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(IDField, pagex.SortAsc)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByIDDESC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(IDField, pagex.SortDesc)
	return qb
}

func (qb *usersRepoQueryBuilder) GroupByID() *usersRepoQueryBuilder {
	qb.QueryBuilder.GroupBy(IDField)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereTenantID(p mysqlx.Predicate, value int64) *usersRepoQueryBuilder {
	qb.QueryBuilder.Where(p, TenantIDField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereTenantIDIn(value []int64) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereIn(TenantIDField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereTenantIDNotIn(value []int64) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereNotIn(TenantIDField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByTenantIDASC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(TenantIDField, pagex.SortAsc)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByTenantIDDESC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(TenantIDField, pagex.SortDesc)
	return qb
}

func (qb *usersRepoQueryBuilder) GroupByTenantID() *usersRepoQueryBuilder {
	qb.QueryBuilder.GroupBy(TenantIDField)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereLoginTime(p mysqlx.Predicate, value time.Time) *usersRepoQueryBuilder {
	qb.QueryBuilder.Where(p, LoginTimeField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereLoginTimeIn(value []time.Time) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereIn(LoginTimeField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereLoginTimeNotIn(value []time.Time) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereNotIn(LoginTimeField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByLoginTimeASC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(LoginTimeField, pagex.SortAsc)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByLoginTimeDESC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(LoginTimeField, pagex.SortDesc)
	return qb
}

func (qb *usersRepoQueryBuilder) GroupByLoginTime() *usersRepoQueryBuilder {
	qb.QueryBuilder.GroupBy(LoginTimeField)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereCreatedAt(p mysqlx.Predicate, value time.Time) *usersRepoQueryBuilder {
	qb.QueryBuilder.Where(p, CreatedAtField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereCreatedAtIn(value []time.Time) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereIn(CreatedAtField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereCreatedAtNotIn(value []time.Time) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereNotIn(CreatedAtField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByCreatedAtASC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(CreatedAtField, pagex.SortAsc)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByCreatedAtDESC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(CreatedAtField, pagex.SortDesc)
	return qb
}

func (qb *usersRepoQueryBuilder) GroupByCreatedAt() *usersRepoQueryBuilder {
	qb.QueryBuilder.GroupBy(CreatedAtField)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereUpdatedAt(p mysqlx.Predicate, value time.Time) *usersRepoQueryBuilder {
	qb.QueryBuilder.Where(p, UpdatedAtField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereUpdatedAtIn(value []time.Time) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereIn(UpdatedAtField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereUpdatedAtNotIn(value []time.Time) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereNotIn(UpdatedAtField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByUpdatedAtASC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(UpdatedAtField, pagex.SortAsc)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByUpdatedAtDESC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(UpdatedAtField, pagex.SortDesc)
	return qb
}

func (qb *usersRepoQueryBuilder) GroupByUpdatedAt() *usersRepoQueryBuilder {
	qb.QueryBuilder.GroupBy(UpdatedAtField)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereDeletedAt(p mysqlx.Predicate, value time.Time) *usersRepoQueryBuilder {
	qb.QueryBuilder.Where(p, DeletedAtField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereDeletedAtIn(value []time.Time) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereIn(DeletedAtField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereDeletedAtNotIn(value []time.Time) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereNotIn(DeletedAtField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByDeletedAtASC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(DeletedAtField, pagex.SortAsc)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByDeletedAtDESC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(DeletedAtField, pagex.SortDesc)
	return qb
}

func (qb *usersRepoQueryBuilder) GroupByDeletedAt() *usersRepoQueryBuilder {
	qb.QueryBuilder.GroupBy(DeletedAtField)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereAccount(p mysqlx.Predicate, value string) *usersRepoQueryBuilder {
	qb.QueryBuilder.Where(p, AccountField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereAccountIn(value []string) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereIn(AccountField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereAccountNotIn(value []string) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereNotIn(AccountField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByAccountASC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(AccountField, pagex.SortAsc)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByAccountDESC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(AccountField, pagex.SortDesc)
	return qb
}

func (qb *usersRepoQueryBuilder) GroupByAccount() *usersRepoQueryBuilder {
	qb.QueryBuilder.GroupBy(AccountField)
	return qb
}

func (qb *usersRepoQueryBuilder) WherePassword(p mysqlx.Predicate, value string) *usersRepoQueryBuilder {
	qb.QueryBuilder.Where(p, PasswordField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WherePasswordIn(value []string) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereIn(PasswordField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WherePasswordNotIn(value []string) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereNotIn(PasswordField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByPasswordASC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(PasswordField, pagex.SortAsc)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByPasswordDESC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(PasswordField, pagex.SortDesc)
	return qb
}

func (qb *usersRepoQueryBuilder) GroupByPassword() *usersRepoQueryBuilder {
	qb.QueryBuilder.GroupBy(PasswordField)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereLoginIP(p mysqlx.Predicate, value string) *usersRepoQueryBuilder {
	qb.QueryBuilder.Where(p, LoginIPField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereLoginIPIn(value []string) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereIn(LoginIPField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) WhereLoginIPNotIn(value []string) *usersRepoQueryBuilder {
	qb.QueryBuilder.WhereNotIn(LoginIPField, value)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByLoginIPASC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(LoginIPField, pagex.SortAsc)
	return qb
}

func (qb *usersRepoQueryBuilder) OrderByLoginIPDESC() *usersRepoQueryBuilder {
	qb.QueryBuilder.OrderBy(LoginIPField, pagex.SortDesc)
	return qb
}

func (qb *usersRepoQueryBuilder) GroupByLoginIP() *usersRepoQueryBuilder {
	qb.QueryBuilder.GroupBy(LoginIPField)
	return qb
}
