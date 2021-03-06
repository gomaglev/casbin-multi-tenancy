package entity

import (
	"context"

	"gin-casbin/internal/app/schema"
	"gin-casbin/pkg/util"

	"github.com/jinzhu/gorm"
)

// GetTenantDB 获取Tenant存储
func GetTenantDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(Tenant))
}

// SchemaTenant 租户对象
type SchemaTenant schema.Tenant

// ToTenant 转换为实体
func (a SchemaTenant) ToTenant() *Tenant {
	item := new(Tenant)
	util.StructMapToStruct(a, item)
	return item
}

// Tenant 租户实体
type Tenant struct {
	Model
	Name              string  `gorm:"column:name;size:200;default:'';not null;"`     // 租户名称
	URL               string  `gorm:"column:url;size:256;default:'';not null;"`      // 租户URL
	LogoURL           string  `gorm:"column:logo_url;size:256;default:'';not null;"` // 租户LOGO URL
	Timezone          string  `gorm:"column:timezone;size:50;default:'';not null;"`  // 时区
	Language          string  `gorm:"column:language;size:50;default:'';not null;"`  // 语言
	Theme             string  `gorm:"column:theme;size:50;default:'';not null;"`     // 默认主题
	Phone             string  `gorm:"size:50;default:'';"`
	Description       *string `gorm:"column:description;"`               // 描述
	Details           *string `gorm:"column:details;"`                   // 详细
	MaxQrQty          int64   `gorm:"column:max_qr_qty;default:100000;"` // 已购买QR数量
	UsedQrQty         int64   `gorm:"column:used_qr_qty;"`               // 已使用QR数量
	TotalOrderQty     int64   `gorm:"column:total_order_qty;"`           // 已提交订单数
	ProcessedOrderQty int64   `gorm:"column:processed_order_qty;"`       // 已处理QR数量
	Status            int     `gorm:"index;default:0;not null;"`         // 状态(1:启用 2:停用)
}

// TableName 表名
func (a Tenant) TableName() string {
	return a.Model.TableName("tenant")
}

// ToSchemaTenant 转换为demo对象
func (a Tenant) ToSchemaTenant() *schema.Tenant {
	item := new(schema.Tenant)
	util.StructMapToStruct(a, item)
	return item
}

// Tenants 租户实体列表
type Tenants []*Tenant

// ToSchemaTenants 转换为对象列表
func (a Tenants) ToSchemaTenants() schema.Tenants {
	list := make(schema.Tenants, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaTenant()
	}
	return list
}
