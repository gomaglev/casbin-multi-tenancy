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
)

var _ bll.ITenant = (*Tenant)(nil)

// TenantSet 注入Tenant
var TenantSet = wire.NewSet(wire.Struct(new(Tenant), "*"), wire.Bind(new(bll.ITenant), new(*Tenant)))

// Tenant 租户
type Tenant struct {
	Enforcer                 *casbin.SyncedEnforcer
	AddressModel             model.IAddress
	TransModel               model.ITrans
	TenantModel              model.ITenant
	UserTenantModel          model.IUserTenant
	UserModel                model.IUser
	UserRoleModel            model.IUserRole
	TenantAdministratorModel model.ITenantAdministrator
	TenantAddressModel       model.ITenantAddress
}

// Query 查询数据
func (a *Tenant) Query(ctx context.Context, params schema.TenantQueryParam, opts ...schema.TenantQueryOptions) (*schema.TenantQueryResult, error) {
	return a.TenantModel.Query(ctx, params, opts...)
}

// Get 查询指定数据
func (a *Tenant) Get(ctx context.Context, id string, opts ...schema.TenantGetOptions) (*schema.Tenant, error) {
	item, err := a.TenantModel.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

// Create 创建数据
func (a *Tenant) Create(ctx context.Context, item schema.Tenant) (*schema.IDResult, error) {
	err := a.checkUserName(ctx, *item.Administrator)
	if err != nil {
		return nil, err
	}

	// tenant id
	tenantID := iutil.NewID()

	// user id / administrator id
	userID := iutil.NewID()

	err = ExecTrans(ctx, a.TransModel, func(ctx context.Context) error {
		// create user role
		userRoles := schema.UserRoles{}
		userRoles = append(userRoles, &schema.UserRole{
			ID:      iutil.NewID(),
			UserID:  userID,
			RoleID:  config.C.TenantOwnerRole.ID,
			Creator: userID,
		})

		for _, urItem := range userRoles {
			if err := a.UserRoleModel.Create(ctx, *urItem); err != nil {
				return err
			}
		}

		userTenant := schema.UserTenant{
			ID:       iutil.NewID(),
			UserID:   userID,
			TenantID: tenantID,
			Creator:  userID,
		}

		// create user tenant
		if err := a.UserTenantModel.Create(ctx, userTenant); err != nil {
			return err
		}

		// create user
		item.Administrator.ID = userID
		item.Administrator.Creator = userID
		item.Administrator.TenantID = tenantID
		item.Administrator.UserRoles = userRoles
		item.Administrator.Password = util.SHA1HashString(item.Administrator.Password)
		if err := a.UserModel.Create(ctx, *item.Administrator); err != nil {
			return err
		}

		// create tenant administrator
		if err := a.TenantAdministratorModel.Create(
			ctx, schema.TenantAdministrator{
				ID:       iutil.NewID(),
				TenantID: tenantID,
				UserID:   userID,
				Creator:  userID,
			}); err != nil {
			return err
		}

		// create tenant address
		addressID := iutil.NewID()
		if err := a.TenantAddressModel.Create(
			ctx, schema.TenantAddress{
				ID:        iutil.NewID(),
				TenantID:  tenantID,
				AddressID: addressID,
				Creator:   userID,
			}); err != nil {
			return err
		}

		// create address
		item.Address.ID = addressID
		item.Address.Creator = userID
		if err := a.AddressModel.Create(ctx, item.Address); err != nil {
			return err
		}

		// create tenant
		item.ID = tenantID
		if err := a.TenantModel.Create(ctx, item); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	LoadCasbinPolicy(ctx, a.Enforcer)
	return schema.NewIDResult(tenantID), nil
}

// Update 更新数据
func (a *Tenant) Update(ctx context.Context, id string, item schema.Tenant) error {
	oldItem, err := a.TenantModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}
	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt

	return a.TenantModel.Update(ctx, id, item)
}

// Delete 删除数据
func (a *Tenant) Delete(ctx context.Context, id string) error {
	oldItem, err := a.TenantModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.TenantModel.Delete(ctx, id)
}

func (a *Tenant) checkUserName(ctx context.Context, item schema.User) error {
	if item.UserName == schema.GetRootUser().UserName {
		return errors.New400Response("ErrIllegalUserName")
	}

	result, err := a.UserModel.Query(ctx,
		schema.UserQueryParam{
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

// UpdateStatus 更新状态
func (a *Tenant) UpdateStatus(ctx context.Context, id string, status int) error {
	oldItem, err := a.TenantModel.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}
	oldItem.Status = status

	err = a.TenantModel.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}
	return nil
}
