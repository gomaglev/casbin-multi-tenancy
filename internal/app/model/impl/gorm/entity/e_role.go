package entity

import (
	"context"

	"gin-casbin/internal/app/schema"
	"gin-casbin/pkg/util"

	"github.com/jinzhu/gorm"
)

// GetRoleDB 获取角色存储
func GetRoleDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(Role))
}

// SchemaRole 角色对象
type SchemaRole schema.Role

// ToRole 转换为角色实体
func (a SchemaRole) ToRole() *Role {
	item := new(Role)
	util.StructMapToStruct(a, item)
	return item
}

// Role 角色实体
type Role struct {
	Model
	Name        string  `gorm:"column:name;size:100;index;default:'';not null;"` // 角色名称
	Sequence    int     `gorm:"column:sequence;index;default:0;not null;"`       // 排序值
	Description *string `gorm:"column:description;size:1024;"`                   // 备注
	Status      int     `gorm:"column:status;index;default:0;not null;"`         // 状态(1:启用 2:禁用)
	Type        int     `gorm:"column:type;index;default:0"`                     // user, owner can manage:0, only root can see, owner:9
}

// TableName 表名
func (a Role) TableName() string {
	return a.Model.TableName("role")
}

// ToSchemaRole 转换为角色对象
func (a Role) ToSchemaRole() *schema.Role {
	item := new(schema.Role)
	util.StructMapToStruct(a, item)
	return item
}

// Roles 角色实体列表
type Roles []*Role

// ToSchemaRoles 转换为角色对象列表
func (a Roles) ToSchemaRoles() []*schema.Role {
	list := make([]*schema.Role, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRole()
	}
	return list
}
