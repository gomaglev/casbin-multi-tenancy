package api

import (
	"gin-casbin/internal/app/bll"
	"gin-casbin/internal/app/config"
	"gin-casbin/internal/app/ginplus"
	"gin-casbin/internal/app/schema"
	"gin-casbin/pkg/errors"
	"gin-casbin/pkg/file"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// TenantSet
var TenantSet = wire.NewSet(wire.Struct(new(Tenant), "*"))

// Tenant
type Tenant struct {
	TenantBll     bll.ITenant
	FileUploadBll bll.IFileUpload
}

// Query
func (a *Tenant) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.TenantQueryParam
	if err := ginplus.ParseQuery(c, &params); err != nil {
		ginplus.ResError(c, err)
		return
	}

	params.Pagination = true
	result, err := a.TenantBll.Query(ctx, params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResPage(c, result.Data, result.PageResult)
}

// Get
func (a *Tenant) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.TenantBll.Get(ctx, c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item)
}

// Create
func (a *Tenant) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Tenant
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	if item.Administrator.Password == "" {
		ginplus.ResError(c, errors.New400Response("ErrPasswordCantBeBlank"))
		return
	}

	result, err := a.TenantBll.Create(ctx, item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, result)
}

// Update
func (a *Tenant) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.Tenant
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := a.TenantBll.Update(ctx, c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Delete
func (a *Tenant) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.TenantBll.Delete(ctx, c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Enable
func (a *Tenant) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.TenantBll.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Disable
func (a *Tenant) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.TenantBll.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Upload
func (a *Tenant) Upload(c *gin.Context) {
	ctx := c.Request.Context()

	tenant, err := a.TenantBll.Get(ctx, c.Param("id"))
	if tenant == nil {
		ginplus.ResError(c, errors.ErrNotFound)
		return
	}

	fileName, fileURL, originalName, fileSize, err := file.Upload(
		c, "tenants", config.C.Upload.MaxFileSize/10, "file", c.Query("from"),
	)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	result, err := a.FileUploadBll.Create(ctx, schema.FileUpload{
		Filename:     fileName,
		FileURL:      fileURL,
		OriginalName: originalName,
		FileSize:     fileSize,
		UserID:       tenant.ID,
		Creator:      ginplus.GetUserID(c),
	})

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	tenant.LogoURL = fileURL
	err = a.TenantBll.Update(ctx, tenant.ID, *tenant)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, result)
}
