package schema

import (
	"gin-casbin/pkg/util"
)

// UserTenant 用户租户对象
type UserTenant struct {
	ID       string `json:"id"`        // 唯一标识
	UserID   string `json:"user_id"`   // 用户ID
	TenantID string `json:"tenant_id"` // 租户ID
	Creator  string `json:"creator"`   // 创建者
}

func (a *UserTenant) String() string {
	return util.JSONMarshalToString(a)
}

// UserTenantQueryParam 查询条件
type UserTenantQueryParam struct {
	UserID  string
	UserIDs []string
	PaginationParam
}

// UserTenantQueryOptions 查询可选参数项
type UserTenantQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// UserTenantGetOptions Get查询可选参数项
type UserTenantGetOptions struct {
}

// UserTenantQueryResult 查询结果
type UserTenantQueryResult struct {
	Data       UserTenants
	PageResult *PaginationResult
}

// UserTenants 用户租户列表
type UserTenants []*UserTenant

// ToTenantIDs 转换为租户ID列表
func (a UserTenants) ToTenantIDs() []string {
	list := make([]string, len(a))
	for i, item := range a {
		list[i] = item.TenantID
	}
	return list
}
