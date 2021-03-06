package api

import (
	"gin-casbin/internal/app/bll"
	"gin-casbin/internal/app/config"
	"gin-casbin/internal/app/ginplus"
	"gin-casbin/internal/app/schema"
	"gin-casbin/pkg/errors"
	"gin-casbin/pkg/logger"

	"github.com/LyricTian/captcha"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// LoginSet
var LoginSet = wire.NewSet(wire.Struct(new(Login), "*"))

// Login
type Login struct {
	LoginBll  bll.ILogin
	QrCodeBll bll.IQrCode
}

// GetCaptcha
func (a *Login) GetCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.LoginBll.GetCaptcha(ctx, config.C.Captcha.Length)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// ResCaptcha
func (a *Login) ResCaptcha(c *gin.Context) {
	ctx := c.Request.Context()
	captchaID := c.Query("id")
	if captchaID == "" {
		ginplus.ResError(c, errors.New400Response("ErrCaptchaIDRequired"))
		return
	}

	if c.Query("reload") != "" {
		if !captcha.Reload(captchaID) {
			ginplus.ResError(c, errors.New400Response("ErrCaptchaIDNotFound"))
			return
		}
	}

	cfg := config.C.Captcha
	err := a.LoginBll.ResCaptcha(ctx, c.Writer, captchaID, cfg.Width, cfg.Height)
	if err != nil {
		ginplus.ResError(c, err)
	}
}

// Login
func (a *Login) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.LoginParam
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if !captcha.VerifyString(item.CaptchaID, item.CaptchaCode) {
		ginplus.ResError(c, errors.New400Response("ErrInvalidCaptchaCode"))
		return
	}

	user, err := a.LoginBll.Verify(ctx, item.UserName, item.Password, c.Request.Referer())
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	userID := user.ID
	tenantID := user.TenantID
	isAdmin := user.IsAdmin

	ginplus.SetUserID(c, userID)
	ginplus.SetTenantID(c, tenantID)
	ginplus.SetIsAdmin(c, isAdmin)

	ctx = logger.NewUserIDContext(ctx, userID, tenantID)
	tokenInfo, err := a.LoginBll.GenerateToken(ctx, userID, tenantID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	logger.StartSpan(ctx, logger.SetSpanTitle("User Login"), logger.SetSpanFuncName("Login")).Infof("登入系统")
	ginplus.ResSuccess(c, tokenInfo)
}

// GetAccessToken Oauth login
func (a *Login) GetAccessToken(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.LoginParam
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	user, err := a.LoginBll.Verify(ctx, item.UserName, item.Password, c.Request.Referer())
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	userID := user.ID
	tenantID := user.TenantID

	ginplus.SetUserID(c, userID)
	ginplus.SetTenantID(c, tenantID)

	ctx = logger.NewUserIDContext(ctx, userID, tenantID)
	tokenInfo, err := a.LoginBll.GenerateToken(ctx, userID, tenantID)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	logger.StartSpan(ctx, logger.SetSpanTitle("User Login"), logger.SetSpanFuncName("Oauth")).Infof("登入系统")
	ginplus.ResSuccess(c, tokenInfo)
}

// Logout
func (a *Login) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	userID := ginplus.GetUserID(c)
	if userID != "" {
		err := a.LoginBll.DestroyToken(ctx, ginplus.GetToken(c))
		if err != nil {
			logger.Errorf(ctx, err.Error())
		}
		logger.StartSpan(ctx, logger.SetSpanTitle("User Logout"), logger.SetSpanFuncName("Logout")).Infof("登出系统")
	}
	ginplus.ResOK(c)
}

// RefreshToken
func (a *Login) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()
	tokenInfo, err := a.LoginBll.GenerateToken(ctx, ginplus.GetUserID(c), ginplus.GetTenantID(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, tokenInfo)
}

// GetUserInfo
func (a *Login) GetUserInfo(c *gin.Context) {
	ctx := c.Request.Context()
	info, err := a.LoginBll.GetLoginInfo(ctx, ginplus.GetUserID(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, info)
}

// QueryUserMenuTree
func (a *Login) QueryUserMenuTree(c *gin.Context) {
	ctx := c.Request.Context()
	menus, err := a.LoginBll.QueryUserMenuTree(ctx, ginplus.GetUserID(c))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResList(c, menus)
}

// UpdatePassword
func (a *Login) UpdatePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.UpdatePasswordParam
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := a.LoginBll.UpdatePassword(ctx, ginplus.GetUserID(c), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// ForgetPassword
func (a *Login) ForgetPassword(c *gin.Context) {
	ctx := c.Request.Context()
	email := c.Query("email")

	err := a.LoginBll.SendResetPasswordMail(ctx, email)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}
