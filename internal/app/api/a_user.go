package api

import (
	"strings"

	"gin-casbin/internal/app/bll"
	"gin-casbin/internal/app/config"
	"gin-casbin/internal/app/ginplus"
	"gin-casbin/internal/app/schema"
	"gin-casbin/pkg/errors"
	"gin-casbin/pkg/file"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

// UserSet
var UserSet = wire.NewSet(wire.Struct(new(User), "*"))

// User
type User struct {
	UserBll       bll.IUser
	FileUploadBll bll.IFileUpload
}

// Query
func (a *User) Query(c *gin.Context) {
	ctx := c.Request.Context()
	var params schema.UserQueryParam
	if err := ginplus.ParseQuery(c, &params); err != nil {
		ginplus.ResError(c, err)
		return
	}
	root := schema.GetRootUser()
	if v := c.Query("roleIDs"); v != "" {
		if params.UserName == root.UserName || params.TenantID == "" {
			params.RoleIDs = strings.Split(v, ",")
		} else {
			roleIDs := strings.Split(v, ",")
			params.RoleIDs = []string{}
			for _, roleID := range roleIDs {
				if roleID != config.C.TenantOwnerRole.ID {
					params.RoleIDs = append(params.RoleIDs, roleID)
				}
			}
		}
	}

	params.Pagination = true
	result, err := a.UserBll.QueryShow(ctx, params)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResPage(c, result.Data, result.PageResult)
}

// Get
func (a *User) Get(c *gin.Context) {
	ctx := c.Request.Context()
	item, err := a.UserBll.Get(ctx, c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, item.CleanSecure())
}

// Create
func (a *User) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	} else if item.Password == "" {
		ginplus.ResError(c, errors.ErrPasswordCantBeBlank)
		return
	}

	item.Creator = ginplus.GetUserID(c)
	tenantID := ginplus.GetTenantID(c)
	if tenantID != "" {
		item.TenantID = tenantID
	}
	result, err := a.UserBll.Create(ctx, item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResSuccess(c, result)
}

// Upload
func (a *User) Upload(c *gin.Context) {
	ctx := c.Request.Context()

	user, err := a.UserBll.Get(ctx, c.Param("id"))
	if user == nil {
		ginplus.ResError(c, errors.ErrNotFound)
		return
	}

	fileName, fileURL, originalName, fileSize, err := file.Upload(
		c, "users", config.C.Upload.MaxFileSize/10, "file", c.Query("from"),
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
		UserID:       user.ID,
		Creator:      ginplus.GetUserID(c),
	})

	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	user.PhotoURL = fileURL
	user.Password = ""
	err = a.UserBll.UpdateProfile(ctx, user.ID, *user)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}

	ginplus.ResSuccess(c, result)
}

// Update
func (a *User) UpdateProfile(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := a.UserBll.UpdateProfile(ctx, c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Update
func (a *User) Update(c *gin.Context) {
	ctx := c.Request.Context()
	var item schema.User
	if err := ginplus.ParseJSON(c, &item); err != nil {
		ginplus.ResError(c, err)
		return
	}

	err := a.UserBll.Update(ctx, c.Param("id"), item)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Delete
func (a *User) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	if ginplus.GetUserID(c) == c.Param("id") {
		ginplus.ResError(c, errors.ErrNotAllowDelete)
		return
	}

	err := a.UserBll.Delete(ctx, c.Param("id"))
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Enable
func (a *User) Enable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserBll.UpdateStatus(ctx, c.Param("id"), 1)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}

// Disable
func (a *User) Disable(c *gin.Context) {
	ctx := c.Request.Context()
	err := a.UserBll.UpdateStatus(ctx, c.Param("id"), 2)
	if err != nil {
		ginplus.ResError(c, err)
		return
	}
	ginplus.ResOK(c)
}
