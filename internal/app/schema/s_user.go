package schema

import (
	"context"
	"time"

	"gin-casbin/internal/app/config"
	"gin-casbin/pkg/util"
)

// GetRootUser 获取root用户
func GetRootUser() *User {
	user := config.C.Root
	return &User{
		ID:       user.UserName,
		UserName: user.UserName,
		RealName: user.RealName,
		Password: util.MD5HashString(user.Password),
	}
}

// CheckIsRootUser 检查是否是root用户
func CheckIsRootUser(ctx context.Context, userID string) bool {
	return GetRootUser().ID == userID
}

// User 用户对象
type User struct {
	ID          string    `json:"id"`                           // 唯一标识
	UserName    string    `json:"user_name" binding:"required"` // 用户名
	RealName    string    `json:"real_name" binding:"required"` // 真实姓名
	Title       string    `json:"title"`                        // 职务
	Password    string    `json:"password"`                     // 密码
	Phone       string    `json:"phone"`                        // 手机号
	Email       string    `json:"email"`                        // 邮箱
	Website     string    `json:"website"`                      // 网站
	PhotoURL    string    `json:"photo_url"`                    // 头像
	Status      int       `json:"status"`                       // 用户状态(1:启用 2:停用)
	Timezone    string    `json:"timezone"`                     // 时区
	Language    string    `json:"language"`                     // 语言
	Theme       string    `json:"theme"`                        // 默认主题
	Description string    `json:"description"`                  // 描述
	Details     string    `json:"details"`                      // 详细
	Creator     string    `json:"creator"`                      // 创建者
	CreatedAt   time.Time `json:"created_at"`                   // 创建时间
	UserRoles   UserRoles `json:"user_roles"`                   // 角色授权
	TenantID    string    `json:"tenant_id"`                    // 租户ID
	IsAdmin     bool      `json:"is_admin"`                     // 是否管理者（临时）
}

func (a *User) String() string {
	return util.JSONMarshalToString(a)
}

// CleanSecure 清理安全数据
func (a *User) CleanSecure() *User {
	a.Password = ""
	return a
}

// UserQueryParam 查询条件
type UserQueryParam struct {
	PaginationParam
	UserName   string   `form:"userName"`   // 用户名
	QueryValue string   `form:"queryValue"` // 模糊查询
	Email      string   `form:"email"`      // EMAIL
	Status     int      `form:"status"`     // 用户状态(1:启用 2:停用)
	RoleIDs    []string `form:"-"`          // 角色ID列表
	TenantID   string   `form:"tenant_id"`  // 租户ID
}

// UserQueryOptions 查询可选参数项
type UserQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// UserQueryResult 查询结果
type UserQueryResult struct {
	Data       Users
	PageResult *PaginationResult
}

// ToShowResult 转换为显示结果
func (a UserQueryResult) ToShowResult(
	mUserRoles map[string]UserRoles,
	mRoles map[string]*Role,
	mUserTenant map[string]UserTenant,
	mTenant map[string]*Tenant) *UserShowQueryResult {
	return &UserShowQueryResult{
		PageResult: a.PageResult,
		Data:       a.Data.ToUserShows(mUserRoles, mRoles, mUserTenant, mTenant),
	}
}

// Users 用户对象列表
type Users []*User

// ToIDs 转换为唯一标识列表
func (a Users) ToIDs() []string {
	idList := make([]string, len(a))
	for i, item := range a {
		idList[i] = item.ID
	}
	return idList
}

// ToUserShows 转换为用户显示列表
func (a Users) ToUserShows(
	mUserRoles map[string]UserRoles,
	mRoles map[string]*Role,
	mUserTenant map[string]UserTenant,
	mTenant map[string]*Tenant) UserShows {
	list := make(UserShows, len(a))
	for i, item := range a {
		showItem := new(UserShow)
		util.StructMapToStruct(item, showItem)
		for _, roleID := range mUserRoles[item.ID].ToRoleIDs() {
			if v, ok := mRoles[roleID]; ok {
				showItem.Roles = append(showItem.Roles, v)
			}
		}
		showItem.Tenant = mTenant[mUserTenant[item.ID].TenantID]
		list[i] = showItem
	}

	return list
}

// ----------------------------------------UserRole--------------------------------------

// UserRole 用户角色
type UserRole struct {
	ID      string `json:"id"`      // 唯一标识
	UserID  string `json:"user_id"` // 用户ID
	RoleID  string `json:"role_id"` // 角色ID
	Creator string `json:"creator"`
}

// UserRoleQueryParam 查询条件
type UserRoleQueryParam struct {
	PaginationParam
	UserID  string   // 用户ID
	UserIDs []string // 用户ID列表
}

// UserRoleQueryOptions 查询可选参数项
type UserRoleQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// UserRoleQueryResult 查询结果
type UserRoleQueryResult struct {
	Data       UserRoles
	PageResult *PaginationResult
}

// UserRoles 角色菜单列表
type UserRoles []*UserRole

// ToMap 转换为map
func (a UserRoles) ToMap() map[string]*UserRole {
	m := make(map[string]*UserRole)
	for _, item := range a {
		m[item.RoleID] = item
	}
	return m
}

// ToRoleIDs 转换为角色ID列表
func (a UserRoles) ToRoleIDs() []string {
	list := make([]string, len(a))
	for i, item := range a {
		list[i] = item.RoleID
	}
	return list
}

// ToUserIDMap 转换为用户ID映射
func (a UserRoles) ToUserIDMap() map[string]UserRoles {
	m := make(map[string]UserRoles)
	for _, item := range a {
		m[item.UserID] = append(m[item.UserID], item)
	}
	return m
}

// ----------------------------------------UserShow--------------------------------------

// UserShow 用户显示项
type UserShow struct {
	ID          string  `json:"id"`                           // 唯一标识
	UserName    string  `json:"user_name" binding:"required"` // 用户名
	RealName    string  `json:"real_name" binding:"required"` // 真实姓名
	Title       string  `json:"title"`                        // 职务
	Phone       string  `json:"phone"`                        // 手机号
	Email       string  `json:"email"`                        // 邮箱
	Website     string  `json:"website"`                      // 网站
	PhotoURL    string  `json:"photo_url"`                    // 头像
	Status      int     `json:"status"`                       // 用户状态(1:启用 2:停用)
	Timezone    string  `json:"timezone"`                     // 时区
	Language    string  `json:"language"`                     // 语言
	Theme       string  `json:"theme"`                        // 默认主题
	Description string  `json:"description"`                  // 描述
	Details     string  `json:"details"`                      // 详细
	Roles       []*Role `json:"roles"`                        // 授权角色列表
	Tenant      *Tenant `json:"tenant"`                       // 租户
}

// UserShows 用户显示项列表
type UserShows []*UserShow

// UserShowQueryResult 用户显示项查询结果
type UserShowQueryResult struct {
	Data       UserShows
	PageResult *PaginationResult
}

// ----------------------------------------UserTenant--------------------------------------
// UserTenant 用户租户
