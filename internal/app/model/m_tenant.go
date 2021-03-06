package model

import (
	"context"

	"gin-casbin/internal/app/schema"
)

// ITenant 租户存储接口s
type ITenant interface {
	// 查询数据
	Query(ctx context.Context, params schema.TenantQueryParam, opts ...schema.TenantQueryOptions) (*schema.TenantQueryResult, error)
	// 查询指定数据
	Get(ctx context.Context, id string, opts ...schema.TenantGetOptions) (*schema.Tenant, error)
	// 创建数据
	Create(ctx context.Context, item schema.Tenant) error
	// 更新数据
	Update(ctx context.Context, id string, item schema.Tenant) error
	// 更新数量
	UpdateCount(ctx context.Context, id string, item schema.Tenant) (int64, error)
	// 删除数据
	Delete(ctx context.Context, id string) error
	// 更新状态
	UpdateStatus(ctx context.Context, id string, status int) error
}
