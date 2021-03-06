package schema

// LoginParam 登录参数
type LoginParam struct {
	UserName    string `json:"user_name" binding:"required"`    // 用户名
	Password    string `json:"password" binding:"required"`     // 密码(md5加密)
	CaptchaID   string `json:"captcha_id" binding:"required"`   // 验证码ID
	CaptchaCode string `json:"captcha_code" binding:"required"` // 验证码
}

// UserLoginInfo 用户登录信息
type UserLoginInfo struct {
	UserID      string `json:"user_id"`   // 用户ID
	UserName    string `json:"user_name"` // 用户名
	RealName    string `json:"real_name"` // 真实姓名
	Email       string `json:"email"`
	PhotoURL    string `json:"photo_url"`
	Phone       string `json:"phone"`
	DisplayName string `json:"display_name"`
	Roles       Roles  `json:"roles"`  // 角色列表
	Tenant      Tenant `json:"tenant"` // 租户
	TenantID    string `json:"tenant_id"`
	Title       string `json:"title"`       // 职务
	Password    string `json:"password"`    // 密码
	Website     string `json:"website"`     // 网站
	Status      int    `json:"status"`      // 用户状态(1:启用 2:停用)
	Timezone    string `json:"timezone"`    // 时区
	Language    string `json:"language"`    // 语言
	Theme       string `json:"theme"`       // 默认主题
	Description string `json:"description"` // 描述
	Details     string `json:"details"`     // 详细
	IsAdmin     bool   `json:"is_admin"`
}

// UpdatePasswordParam 更新密码请求参数
type UpdatePasswordParam struct {
	OldPassword string `json:"old_password" binding:"required"` // 旧密码(md5加密)
	NewPassword string `json:"new_password" binding:"required"` // 新密码(md5加密)
}

// ResetPasswordParam 重置密码请求参数
type ResetPasswordParam struct {
	Email string `json:"email" binding:"required"` // EMAIL
}

// LoginCaptcha 登录验证码
type LoginCaptcha struct {
	CaptchaID string `json:"captcha_id"` // 验证码ID
	Type      string `json:"type"`
}

// LoginTokenInfo 登录令牌信息
type LoginTokenInfo struct {
	AccessToken string `json:"access_token"` // 访问令牌
	TokenType   string `json:"token_type"`   // 令牌类型
	ExpiresAt   int64  `json:"expires_at"`   // 令牌到期时间戳
}
