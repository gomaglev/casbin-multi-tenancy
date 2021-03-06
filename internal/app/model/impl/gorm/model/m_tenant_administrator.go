package model

import (
	"context"

	"gin-casbin/internal/app/model"
	"gin-casbin/internal/app/model/impl/gorm/entity"
	"gin-casbin/internal/app/schema"
	"gin-casbin/pkg/errors"

	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var _ model.ITenantAdministrator = (*TenantAdministrator)(nil)

// TenantAdministratorSet 注入TenantAdministrator
var TenantAdministratorSet = wire.NewSet(wire.Struct(new(TenantAdministrator), "*"), wire.Bind(new(model.ITenantAdministrator), new(*TenantAdministrator)))

// TenantAdministrator 租户主用户存储
type TenantAdministrator struct {
	DB *gorm.DB
}

func (a *TenantAdministrator) getQueryOption(opts ...schema.TenantAdministratorQueryOptions) schema.TenantAdministratorQueryOptions {
	var opt schema.TenantAdministratorQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *TenantAdministrator) Query(ctx context.Context, params schema.TenantAdministratorQueryParam, opts ...schema.TenantAdministratorQueryOptions) (*schema.TenantAdministratorQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetTenantAdministratorDB(ctx, a.DB)

	if v := params.TenantID; v != "" {
		db = db.Where("tenant_id=?", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.TenantAdministrators
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.TenantAdministratorQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaTenantAdministrators(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *TenantAdministrator) Get(ctx context.Context, id string, opts ...schema.TenantAdministratorGetOptions) (*schema.TenantAdministrator, error) {
	db := entity.GetTenantAdministratorDB(ctx, a.DB).Where("id=?", id)
	var item entity.TenantAdministrator
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaTenantAdministrator(), nil
}

// Create 创建数据
func (a *TenantAdministrator) Create(ctx context.Context, item schema.TenantAdministrator) error {
	eitem := entity.SchemaTenantAdministrator(item).ToTenantAdministrator()
	result := entity.GetTenantAdministratorDB(ctx, a.DB).Create(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *TenantAdministrator) Update(ctx context.Context, id string, item schema.TenantAdministrator) error {
	eitem := entity.SchemaTenantAdministrator(item).ToTenantAdministrator()
	result := entity.GetTenantAdministratorDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *TenantAdministrator) Delete(ctx context.Context, id string) error {
	result := entity.GetTenantAdministratorDB(ctx, a.DB).Where("id=?", id).Delete(entity.TenantAdministrator{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
