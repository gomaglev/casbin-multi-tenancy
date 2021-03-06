package model

import (
	"context"
	"time"

	"gin-casbin/internal/app/model"
	"gin-casbin/internal/app/model/impl/gorm/entity"
	"gin-casbin/internal/app/schema"
	"gin-casbin/pkg/errors"

	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var _ model.ITenant = (*Tenant)(nil)

// TenantSet 注入Tenant
var TenantSet = wire.NewSet(wire.Struct(new(Tenant), "*"), wire.Bind(new(model.ITenant), new(*Tenant)))

// Tenant 租户存储
type Tenant struct {
	DB *gorm.DB
}

func (a *Tenant) getQueryOption(opts ...schema.TenantQueryOptions) schema.TenantQueryOptions {
	var opt schema.TenantQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

// Query 查询数据
func (a *Tenant) Query(ctx context.Context, params schema.TenantQueryParam, opts ...schema.TenantQueryOptions) (*schema.TenantQueryResult, error) {
	opt := a.getQueryOption(opts...)

	db := entity.GetTenantDB(ctx, a.DB)
	if v := params.QrCodeID; v != "" {
		db = db.Joins(`join qr_code_tenant on qr_code_tenant.tenant_id = tenant.id 
		               and qr_code_tenant.qr_code_id=?`, v)
	}

	if v := params.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	opt.OrderFields = append(opt.OrderFields, schema.NewOrderField("id", schema.OrderByDESC))
	db = db.Order(ParseOrder(opt.OrderFields))

	var list entity.Tenants
	pr, err := WrapPageQuery(ctx, db, params.PaginationParam, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	qr := &schema.TenantQueryResult{
		PageResult: pr,
		Data:       list.ToSchemaTenants(),
	}

	return qr, nil
}

// Get 查询指定数据
func (a *Tenant) Get(ctx context.Context, id string, opts ...schema.TenantGetOptions) (*schema.Tenant, error) {
	db := entity.GetTenantDB(ctx, a.DB).Where("id=?", id)

	var item entity.Tenant
	ok, err := FindOne(ctx, db, &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaTenant(), nil
}

// Create 创建数据
func (a *Tenant) Create(ctx context.Context, item schema.Tenant) error {
	eitem := entity.SchemaTenant(item).ToTenant()
	result := entity.GetTenantDB(ctx, a.DB).Create(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update 更新数据
func (a *Tenant) Update(ctx context.Context, id string, item schema.Tenant) error {
	eitem := entity.SchemaTenant(item).ToTenant()
	result := entity.GetTenantDB(ctx, a.DB).Where("id=?", id).Updates(eitem)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Delete 删除数据
func (a *Tenant) Delete(ctx context.Context, id string) error {
	result := entity.GetTenantDB(ctx, a.DB).Where("id=?", id).Delete(entity.Tenant{})
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateStatus 更新状态
func (a *Tenant) UpdateStatus(ctx context.Context, id string, status int) error {
	result := entity.GetTenantDB(ctx, a.DB).Where("id=?", id).Update("status", status)
	if err := result.Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpdateCount
func (m *Tenant) UpdateCount(ctx context.Context, id string, item schema.Tenant) (int64, error) {
	countMap := make(map[string]interface{})

	if item.MaxQrQty > 0 {
		countMap["max_qr_qty"] = gorm.Expr("max_qr_qty + ?", item.MaxQrQty)
	}
	if item.UsedQrQty > 0 {
		countMap["used_qr_qty"] = gorm.Expr("used_qr_qty + ?", item.UsedQrQty)
	}
	if item.TotalOrderQty > 0 {
		countMap["total_order_qty"] = gorm.Expr("total_order_qty + ?", item.TotalOrderQty)
	}
	if item.ProcessedOrderQty > 0 {
		countMap["processed_order_qty"] = gorm.Expr("processed_order_qty + ?", item.ProcessedOrderQty)
	}

	if len(countMap) > 0 {
		countMap["updated_at"] = time.Now()
		result := entity.GetTenantDB(ctx, m.DB).Where("id=?", id).UpdateColumns(countMap)
		if err := result.Error; err != nil {
			return result.RowsAffected, errors.WithStack(err)
		}
		return result.RowsAffected, nil
	}
	return 0, nil
}
