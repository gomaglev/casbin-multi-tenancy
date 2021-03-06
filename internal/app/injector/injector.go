package injector

import (
	"gin-casbin/pkg/auth"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// InjectorSet
var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

// Injector - engine and Casbin Enforcer
type Injector struct {
	Engine         *gin.Engine
	Auth           auth.Auther
	CasbinEnforcer *casbin.SyncedEnforcer
}
