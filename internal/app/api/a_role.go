package api

import (
	"gin-casbin/internal/app/bll"
	"gin-casbin/internal/app/ginplus"
	"gin-casbin/internal/app/schema"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// RoleSet 注入Role
var RoleSet = wire.NewSet(wire.Struct(new(Role), "*"))

// Role
type Role struct {
	RoleBll bll.IRole
}

// Query
func (a *Role) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RoleQueryParam
	if err := ginplus.ParseQuery(c, &params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	params.Pagination = true
	params.TenantID = ginplus.GetTenantID(c)
	result, err := a.RoleBll.Query(ctx, params, schema.RoleQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResPage(c, result.Data, result.PageResult)
}

// QuerySelect
func (a *Role) QuerySelect(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.RoleQueryParam
	if err := ginplus.ParseQuery(c, &params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	result, err := a.RoleBll.Query(ctx, params, schema.RoleQueryOptions{
		OrderFields: schema.NewOrderFields(schema.NewOrderField("sequence", schema.OrderByDESC)),
	})
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResList(c, result.Data)
}

// Get
func (a *Role) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.RoleBll.Get(ctx, c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// Create
func (a *Role) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Role
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	item.Creator = ginplus.GetUserID(c)
	result, err := a.RoleBll.Create(ctx, item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, result)
}

// Update
func (a *Role) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Role
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := a.RoleBll.Update(ctx, c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Delete
func (a *Role) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleBll.Delete(ctx, c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Enable
func (a *Role) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleBll.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Disable
func (a *Role) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.RoleBll.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}
