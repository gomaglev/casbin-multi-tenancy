package model

import (
	"context"
	"log"

	"gin-casbin/internal/app/model"
	"gin-casbin/internal/app/model/impl/gorm/entity"
	"gin-casbin/internal/app/schema"
	"gin-casbin/pkg/errors"

	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var _ model.IUserTenant = (*UserTenant)(nil)

// UserTenantSet 注入UserTenant
var UserTenantSet = wire.NewSet(wire.Struct(new(UserTenant), "*"), wire.Bind(new(model.IUserTenant), new(*UserTenant)))

// UserTenant 用户租户存储
type UserTenant struct {
	DB *gorm.DB
}

func (a *UserTenant) getQueryOption(opts ...schema.UserTenantQueryOptions) schema.UserTenantQueryOptions {
	var opt schema.UserTenantQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *UserTenant) Query(ctx context.Context, params schema.UserTenantQueryParam, opts ...schema.UserTenantQueryOptions) (*schema.UserTenantQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetUserTenantDB(ctx, a.DB)
	if v := params.UserID; v != "" {
		db = db.Where("user_id=?", v)
	}
	if v := params.UserIDs; len(v) > 0 {
		db = db.Where("user_id IN (?)", v)
	}
	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.UserTenants
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.UserTenantQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaUserTenants(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *UserTenant) Get(ctx context.Context, id string, opts ...schema.UserTenantGetOptions) (*schema.UserTenant, error) {
	db := entity.GetUserTenantDB(ctx, a.DB).Where("id=?", id)
	var item entity.UserTenant
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaUserTenant(), nil
}

// Create 创建数据
func (a *UserTenant) Create(ctx context.Context, item schema.UserTenant) error {
	eitem := entity.SchemaUserTenant(item).ToUserTenant()
	log.Printf("entity.SchemaUserTenant(item).ToUserTenant(): %s", eitem)
	result := entity.GetUserTenantDB(ctx, a.DB).Create(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *UserTenant) Update(ctx context.Context, id string, item schema.UserTenant) error {
	eitem := entity.SchemaUserTenant(item).ToUserTenant()
	result := entity.GetUserTenantDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *UserTenant) Delete(ctx context.Context, id string) error {
	result := entity.GetUserTenantDB(ctx, a.DB).Where("id=?", id).Delete(entity.UserTenant{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
