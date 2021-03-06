package bll

import (
	"context"

	"gin-casbin/internal/app/bll"
	"gin-casbin/internal/app/config"
	"gin-casbin/internal/app/iutil"
	"gin-casbin/internal/app/model"
	"gin-casbin/internal/app/schema"
	"gin-casbin/pkg/errors"
	"gin-casbin/pkg/util"

	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

var _ bll.IUser = (*User)(nil)

// UserSet 注入User
var UserSet = wire.NewSet(wire.Struct(new(User), "*"), wire.Bind(new(bll.IUser), new(*User)))

// User 用户管理
type User struct {
	Enforcer        *casbin.SyncedEnforcer
	TransModel      model.ITrans
	UserModel       model.IUser
	UserRoleModel   model.IUserRole
	RoleModel       model.IRole
	UserTenantModel model.IUserTenant
	TenantModel     model.ITenant
}

// Query 查询数据
func (a *User) Query(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserQueryResult, error) {
	return a.UserModel.Query(ctx, params, opts...)
}

// QueryShow 查询显示项数据
func (a *User) QueryShow(ctx context.Context, params schema.UserQueryParam, opts ...schema.UserQueryOptions) (*schema.UserShowQueryResult, error) {
	users, err := a.UserModel.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	} else if users == nil {
		return nil, nil
	}

	userIDs := users.Data.ToIDs()
	userRoleResult, err := a.UserRoleModel.Query(ctx, schema.UserRoleQueryParam{
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}

	roleResult, err := a.RoleModel.Query(ctx, schema.RoleQueryParam{
		IDs: userRoleResult.Data.ToRoleIDs(),
	})
	if err != nil {
		return nil, err
	}

	userTenantResult, err := a.UserTenantModel.Query(ctx, schema.UserTenantQueryParam{
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}

	tenantIDs := userTenantResult.Data.ToTenantIDs()
	tenantMap := make(map[string]*schema.Tenant)
	if len(tenantIDs) > 0 {
		tenantResult, err := a.TenantModel.Query(ctx, schema.TenantQueryParam{
			IDs: tenantIDs,
		})
		if err != nil {
			return nil, err
		}
		logrus.Printf("tenantResult, %v", tenantResult)
		tenantMap = tenantResult.Data.ToMap()
	}

	return users.ToShowResult(
		userRoleResult.Data.ToUserIDMap(),
		roleResult.Data.ToMap(),
		userTenantResult.Data.ToUserIDMap(),
		tenantMap,
	), nil
}

// Get 查询指定数据
func (a *User) Get(ctx context.Context, id string, opts ...schema.UserQueryOptions) (*schema.User, error) {
	item, err := a.UserModel.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	userRoleResult, err := a.UserRoleModel.Query(ctx, schema.UserRoleQueryParam{
		UserID: id,
	})
	if err != nil {
		return nil, err
	}
	item.UserRoles = userRoleResult.Data

	userTenantResult, err := a.UserTenantModel.Query(ctx, schema.UserTenantQueryParam{
		UserID: id,
	})
	if err != nil {
		return nil, err
	}
	if userTenantResult.Data != nil && len(userTenantResult.Data) > 0 {
		item.TenantID = userTenantResult.Data[0].TenantID
	}

	// first version only use 2 roles: admin & non-admin
	for _, userRole := range item.UserRoles {
		if userRole.RoleID == config.C.TenantOwnerRole.ID {
			item.IsAdmin = true
		}
	}
	return item, nil
}

// Create 创建数据
func (a *User) Create(ctx context.Context, item schema.User) (*schema.IDResult, error) {
	err := a.checkUserName(ctx, item)
	if err != nil {
		return nil, err
	}

	item.Password = util.SHA1HashString(item.Password)
	item.ID = iutil.NewID()
	err = ExecTrans(ctx, a.TransModel, func(ctx context.Context) error {
		for _, urItem := range item.UserRoles {
			urItem.ID = iutil.NewID()
			urItem.UserID = item.ID
			err := a.UserRoleModel.Create(ctx, *urItem)
			if err != nil {
				return err
			}
		}
		logrus.Printf("UserID:%s,TenantID:%s,item.UserRoles:%v", item.ID, item.TenantID, item.UserRoles)
		err := a.UserTenantModel.Create(ctx, schema.UserTenant{
			ID:       iutil.NewID(),
			UserID:   item.ID,
			TenantID: item.TenantID,
			Creator:  item.Creator,
		})
		if err != nil {
			return err
		}

		return a.UserModel.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return schema.NewIDResult(item.ID), nil
}

func (a *User) checkUserName(ctx context.Context, item schema.User) error {
	if item.UserName == schema.GetRootUser().UserName {
		return errors.New400Response("ErrIllegalUserName")
	}
	result, err := a.UserModel.Query(ctx, schema.UserQueryParam{
		PaginationParam: schema.PaginationParam{OnlyCount: true},
		UserName:        item.UserName,
	})
	if err != nil {
		return err
	} else if result.PageResult.Total > 0 {
		return errors.New400Response("ErrDuplicatedUserName")
	}
	return nil
}

// Update 更新数据
func (a *User) Update(ctx context.Context, id string, item schema.User) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.UserName != item.UserName {
		err := a.checkUserName(ctx, item)
		if err != nil {
			return err
		}
	}

	if item.Password != "" {
		item.Password = util.SHA1HashString(item.Password)
	} else {
		item.Password = oldItem.Password
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt
	err = ExecTrans(ctx, a.TransModel, func(ctx context.Context) error {
		addUserRoles, delUserRoles := a.compareUserRoles(ctx, oldItem.UserRoles, item.UserRoles)
		for _, rmitem := range addUserRoles {
			rmitem.ID = iutil.NewID()
			rmitem.UserID = id
			err := a.UserRoleModel.Create(ctx, *rmitem)
			if err != nil {
				return err
			}
		}

		for _, rmitem := range delUserRoles {
			err := a.UserRoleModel.Delete(ctx, rmitem.ID)
			if err != nil {
				return err
			}
		}

		return a.UserModel.Update(ctx, id, item)
	})
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}

// UpdateProfile 更新个人资料
func (a *User) UpdateProfile(ctx context.Context, id string, item schema.User) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.UserName != item.UserName {
		err := a.checkUserName(ctx, item)
		if err != nil {
			return err
		}
	}

	if item.Password != "" {
		item.Password = util.SHA1HashString(item.Password)
	} else {
		item.Password = oldItem.Password
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt
	err = ExecTrans(ctx, a.TransModel, func(ctx context.Context) error {
		return a.UserModel.Update(ctx, id, item)
	})
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}

func (a *User) compareUserRoles(ctx context.Context, oldUserRoles, newUserRoles schema.UserRoles) (addList, delList schema.UserRoles) {
	mOldUserRoles := oldUserRoles.ToMap()
	mNewUserRoles := newUserRoles.ToMap()

	for k, item := range mNewUserRoles {
		if _, ok := mOldUserRoles[k]; ok {
			delete(mOldUserRoles, k)
			continue
		}
		addList = append(addList, item)
	}

	for _, item := range mOldUserRoles {
		delList = append(delList, item)
	}
	return
}

// Delete 删除数据
func (a *User) Delete(ctx context.Context, id string) error {
	oldItem, err := a.UserModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	err = ExecTrans(ctx, a.TransModel, func(ctx context.Context) error {
		err := a.UserRoleModel.DeleteByUserID(ctx, id)
		if err != nil {
			return err
		}

		return a.UserModel.Delete(ctx, id)
	})
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}

// UpdateStatus 更新状态
func (a *User) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.UserModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}
	oldItem.Status = status

	err = a.UserModel.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	LoadCasbinPolicy(ctx, a.Enforcer)
	return nil
}
