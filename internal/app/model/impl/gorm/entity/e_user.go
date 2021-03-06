package entity

import (
	"context"

	"gin-casbin/internal/app/schema"
	"gin-casbin/pkg/util"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// GetUserDB 获取用户存储
func GetUserDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return GetDBWithModel(ctx, defDB, new(User))
}

// SchemaUser 用户对象
type SchemaUser schema.User

// ToUser 转换为用户实体
func (a SchemaUser) ToUser() *User {
	item := new(User)
	util.StructMapToStruct(a, item)
	return item
}

// User 用户实体
type User struct {
	Model
	UserName    string  `gorm:"size:64;index;default:'';not null;"` // 用户名
	RealName    string  `gorm:"size:64;index;default:'';not null;"` // 真实姓名
	Title       string  `gorm:"size:64;index;default:'';not null;"` // 职务
	Password    string  `gorm:"size:255;default:'';not null;"`      // 密码(sha1(md5(明文))加密)
	Email       *string `gorm:"size:255;index;"`                    // 邮箱
	Phone       *string `gorm:"size:50;index;"`                     // 手机号
	Website     *string `gorm:"size:255;index;"`                    // 网站
	PhotoURL    *string `gorm:"size:512;"`                          // 头像
	Status      int     `gorm:"index;default:0;not null;"`          // 状态(1:启用 2:停用)
	Timezone    string  `gorm:"size:50;default:'';not null;"`       // 时区
	Language    string  `gorm:"size:50;default:'';not null;"`       // 语言
	Theme       string  `gorm:"size:50;default:'';not null;"`       // 默认主题
	TenantID    string  `gorm:"size:36;index;default:'';not null;"` // 租户ID
	Tenant      Tenant  // user belongs to tenant
	Roles       []Role  `gorm:"many2many:user_role;"`
	Description *string // 描述
	Details     *string // 详细
}

// TableName 表名
func (a User) TableName() string {
	return a.Model.TableName("user")
}

// ToSchemaUser 转换为用户对象
func (a User) ToSchemaUser() *schema.User {
	item := new(schema.User)
	util.StructMapToStruct(a, item)

	logrus.Printf("from:%v,to:%v", a, item)
	return item
}

// Users 用户实体列表
type Users []*User

// ToSchemaUsers 转换为用户对象列表
func (a Users) ToSchemaUsers() []*schema.User {
	list := make([]*schema.User, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaUser()
	}
	return list
}
