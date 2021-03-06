package router

import (
	"gin-casbin/internal/app/api"
	"gin-casbin/pkg/auth"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var _ IRouter = (*Router)(nil)

// RouterSet
var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

// IRouter
type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

// Router
type Router struct {
	Auth           auth.Auther
	CasbinEnforcer *casbin.SyncedEnforcer
	LoginAPI       *api.Login
	RoleAPI        *api.Role
	UserAPI        *api.User
	TenantAPI      *api.Tenant
	ResourceAPI    *api.Resource
}

// Register
func (a *Router) Register(app *gin.Engine) error {
	a.RegisterAPI(app)
	return nil
}

// Prefixes
func (a *Router) Prefixes() []string {
	return []string{
		"/api/",
	}
}
