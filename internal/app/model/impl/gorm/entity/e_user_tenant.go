package entity

import (
	"context"

	"gin-casbin/internal/app/schema"
	"gin-casbin/pkg/util"

	"github.com/jinzhu/gorm"
)

// GetUserTenantDB 获取UserTenant存储
func GetUserTenantDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(UserTenant))
}

// SchemaUserTenant 用户租户对象
type SchemaUserTenant schema.UserTenant

// ToUserTenant 转换为实体
func (a SchemaUserTenant) ToUserTenant() *UserTenant {
	item := new(UserTenant)
	util.StructMapToStruct(a, item)
	return item
}

// UserTenant 用户租户实体
type UserTenant struct {
	Model
	UserID   string `gorm:"column:user_id;size:36;index;default:'';not null;"`   // 用户ID
	TenantID string `gorm:"column:tenant_id;size:36;index;default:'';not null;"` // 租户ID
}

// TableName 表名
func (a UserTenant) TableName() string {
	return a.Model.TableName("user_tenant")
}

// ToSchemaUserTenant 转换为demo对象
func (a UserTenant) ToSchemaUserTenant() *schema.UserTenant {
	item := new(schema.UserTenant)
	util.StructMapToStruct(a, item)
	return item
}

// UserTenants 用户租户实体列表
type UserTenants []*UserTenant

// ToSchemaUserTenants 转换为对象列表
func (a UserTenants) ToSchemaUserTenants() schema.UserTenants {
	list := make(schema.UserTenants, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUserTenant()
	}
	return list
}
