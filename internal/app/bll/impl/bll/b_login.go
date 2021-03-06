package bll

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"gin-casbin/internal/app/bll"
	"gin-casbin/internal/app/config"
	"gin-casbin/internal/app/model"
	"gin-casbin/internal/app/schema"
	"gin-casbin/pkg/auth"
	"gin-casbin/pkg/errors"
	"gin-casbin/pkg/mail"
	"gin-casbin/pkg/util"

	"github.com/LyricTian/captcha"
	"github.com/google/wire"
	"gopkg.in/gomail.v2"
)

var _ bll.ILogin = (*Login)(nil)

// LoginSet 注入Login
var LoginSet = wire.NewSet(wire.Struct(new(Login), "*"), wire.Bind(new(bll.ILogin), new(*Login)))

// Login 登录管理
type Login struct {
	Auth            auth.Auther
	UserModel       model.IUser
	UserRoleModel   model.IUserRole
	UserTenantModel model.IUserTenant
	TenantModel     model.ITenant
	RoleModel       model.IRole
	RoleMenuModel   model.IRoleMenu
	MenuModel       model.IMenu
	MenuActionModel model.IMenuAction
	Mailer          *mail.Mailer
}

// GetCaptcha 获取图形验证码信息
func (a *Login) GetCaptcha(ctx context.Context, length int) (*schema.LoginCaptcha, error) {
	captchaID := captcha.NewLen(length)
	item := &schema.LoginCaptcha{
		CaptchaID: captchaID,
	}
	return item, nil
}

// ResCaptcha 生成并响应图形验证码
func (a *Login) ResCaptcha(ctx context.Context, w http.ResponseWriter, captchaID string, width, height int) error {
	err := captcha.WriteImage(w, captchaID, width, height)
	if err != nil {
		if err == captcha.ErrNotFound {
			return errors.ErrNotFound
		}
		return errors.WithStack(err)
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "image/png")
	return nil
}

// Verify 登录验证
func (a *Login) Verify(ctx context.Context, userName, password string, referer string) (*schema.User, error) {
	// 检查是否是超级用户
	root := schema.GetRootUser()
	if userName == root.UserName && root.Password == password &&
		!strings.HasSuffix(strings.ToLower(referer), "sessions/signin") {
		root.TenantID = "wetrue"
		return root, nil
	}

	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		UserName: userName,
	})
	if err != nil {
		return nil, err
	} else if len(result.Data) == 0 {
		return nil, errors.ErrInvalidUserName
	}
	item := result.Data[0]
	if item.Password != util.SHA1HashString(password) {
		return nil, errors.ErrInvalidPassword
	} else if item.Status != 1 {
		return nil, errors.ErrUserDisable
	}

	return item, nil
}

// GenerateToken 生成令牌
func (a *Login) GenerateToken(ctx context.Context, userID string, tenantID string) (*schema.LoginTokenInfo, error) {
	tokenInfo, err := a.Auth.GenerateToken(ctx, userID, tenantID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	item := &schema.LoginTokenInfo{
		AccessToken: tokenInfo.GetAccessToken(),
		TokenType:   tokenInfo.GetTokenType(),
		ExpiresAt:   tokenInfo.GetExpiresAt(),
	}
	return item, nil
}

// DestroyToken 销毁令牌
func (a *Login) DestroyToken(ctx context.Context, tokenString string) error {
	err := a.Auth.DestroyToken(ctx, tokenString)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (a *Login) checkAndGetUser(ctx context.Context, userID string) (*schema.User, error) {
	user, err := a.UserModel.Get(ctx, userID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.ErrInvalidUser
	} else if user.Status != 1 {
		return nil, errors.ErrUserDisable
	}
	return user, nil
}

// GetLoginInfo 获取当前用户登录信息
func (a *Login) GetLoginInfo(ctx context.Context, userID string) (*schema.UserLoginInfo, error) {
	if isRoot := schema.CheckIsRootUser(ctx, userID); isRoot {
		root := schema.GetRootUser()
		loginInfo := &schema.UserLoginInfo{
			UserName: root.UserName,
			RealName: root.RealName,
		}
		return loginInfo, nil
	}

	user, err := a.checkAndGetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	info := &schema.UserLoginInfo{
		UserID:      user.ID,
		UserName:    user.UserName,
		RealName:    user.RealName,
		Email:       user.Email,
		PhotoURL:    user.PhotoURL,
		Phone:       user.Phone,
		DisplayName: user.RealName,
		Title:       user.Title,
		Password:    user.Password,
		Website:     user.Website,
		Status:      user.Status,
		Timezone:    user.Timezone,
		Language:    user.Language,
		Theme:       user.Theme,
		Description: user.Description,
		Details:     user.Details,
		IsAdmin:     user.IsAdmin,
	}

	userRoleResult, err := a.UserRoleModel.Query(ctx, schema.UserRoleQueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	if roleIDs := userRoleResult.Data.ToRoleIDs(); len(roleIDs) > 0 {
		roleResult, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
			IDs:    roleIDs,
			Status: 1,
		})
		if err != nil {
			return nil, err
		}
		info.Roles = roleResult.Data
		for _, roleID := range roleIDs {
			if roleID == config.C.TenantOwnerRole.ID {
				info.IsAdmin = true
				break
			}
		}
	}

	tenantResult, err := a.TenantModel.Get(ctx, user.TenantID)
	if err != nil {
		return nil, err
	}
	info.TenantID = user.TenantID
	info.Tenant = *tenantResult

	return info, nil
}

// QueryUserMenuTree 查询当前用户的权限菜单树
func (a *Login) QueryUserMenuTree(ctx context.Context, userID string) (schema.MenuTrees, error) {
	isRoot := schema.CheckIsRootUser(ctx, userID)
	// 如果是root用户，则查询所有显示的菜单树
	if isRoot {
		result, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
			Status: 1,
		}, schema.MenuQueryOptions{
			OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
		})
		if err != nil {
			return nil, err
		}

		menuActionResult, err := a.MenuActionModel.Query(ctx, schema.MenuActionQueryParam{})
		if err != nil {
			return nil, err
		}
		return result.Data.FillMenuAction(menuActionResult.Data.ToMenuIDMap()).ToTree(), nil
	}

	userRoleResult, err := a.UserRoleModel.Query(ctx, schema.UserRoleQueryParam{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	} else if len(userRoleResult.Data) == 0 {
		return nil, errors.ErrNoPerm
	}

	roleMenuResult, err := a.RoleMenuModel.Query(ctx, schema.RoleMenuQueryParam{
		RoleIDs: userRoleResult.Data.ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	} else if len(roleMenuResult.Data) == 0 {
		return nil, errors.ErrNoPerm
	}

	menuResult, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
		IDs:    roleMenuResult.Data.ToMenuIDs(),
		Status: 1,
	})
	if err != nil {
		return nil, err
	} else if len(menuResult.Data) == 0 {
		return nil, errors.ErrNoPerm
	}

	mData := menuResult.Data.ToMap()
	var qIDs []string
	for _, pid := range menuResult.Data.SplitParentIDs() {
		if _, ok := mData[pid]; !ok {
			qIDs = append(qIDs, pid)
		}
	}

	if len(qIDs) > 0 {
		pmenuResult, err := a.MenuModel.Query(ctx, schema.MenuQueryParam{
			IDs: menuResult.Data.SplitParentIDs(),
		})
		if err != nil {
			return nil, err
		}
		menuResult.Data = append(menuResult.Data, pmenuResult.Data...)
	}

	sort.Sort(menuResult.Data)
	menuActionResult, err := a.MenuActionModel.Query(ctx, schema.MenuActionQueryParam{
		IDs: roleMenuResult.Data.ToActionIDs(),
	})
	if err != nil {
		return nil, err
	}
	return menuResult.Data.FillMenuAction(menuActionResult.Data.ToMenuIDMap()).ToTree(), nil
}

// UpdatePassword 更新当前用户登录密码
func (a *Login) UpdatePassword(ctx context.Context, userID string, params schema.UpdatePasswordParam) error {
	if schema.CheckIsRootUser(ctx, userID) {
		return errors.New400Response("can't change root password")
	}

	user, err := a.checkAndGetUser(ctx, userID)
	if err != nil {
		return err
	} else if util.SHA1HashString(params.OldPassword) != user.Password {
		return errors.New400Response("ErrWrongOldPassword")
	}

	params.NewPassword = util.SHA1HashString(params.NewPassword)
	return a.UserModel.UpdatePassword(ctx, userID, params.NewPassword)
}

// SendResetPasswordMail 发送更改密码邮件
func (a *Login) SendResetPasswordMail(ctx context.Context, email string) error {
	user, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		Email:  email,
		Status: 1,
	})
	if err != nil {
		return err
	}
	if user != nil && len(user.Data) > 0 {
		m := gomail.NewMessage()

		m.SetHeader("From", config.C.SMTP.From)
		m.SetHeader("To", user.Data[0].Email)
		m.SetHeader("Subject", "WeTrue - Reset Password")
		url := ""
		for _, user := range user.Data {
			url = url + fmt.Sprintf(
				`<a href="http://2t.ai/reset-password?username=%s">Change for user %s.</a><br>`, user.UserName, user.UserName)
		}

		m.SetBody("text/html", fmt.Sprintf(
			`Please use the link below to change your password. <br>%s`, url))

		a.Mailer.SendChan <- m
	}

	return nil
}
