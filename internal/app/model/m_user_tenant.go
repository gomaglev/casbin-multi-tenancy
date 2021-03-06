package model

import (
	"context"

	"gin-casbin/internal/app/schema"
)

// IUserTenant 用户租户存储接口
type IUserTenant interface {
	// 查询数据
	Query(ctx context.Context, params schema.UserTenantQueryParam, opts ...schema.UserTenantQueryOptions) (*schema.UserTenantQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, id string, opts ...schema.UserTenantGetOptions) (*schema.UserTenant, error)
	// 创建数据
	Create(ctx context.Context, item schema.UserTenant) error
	// 更新数据
	Update(ctx context.Context, id string, item schema.UserTenant) error
	// 删除数据
	Delete(ctx context.Context, id string) error
}
