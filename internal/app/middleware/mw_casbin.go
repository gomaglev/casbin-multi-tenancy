package middleware

import (
	"gin-casbin/internal/app/config"
	"gin-casbin/internal/app/ginplus"
	"gin-casbin/pkg/errors"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// CasbinMiddleware casbin中间件
func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) gin.HandlerFunc {
	cfg := config.C.Casbin
	if !cfg.Enable {
		return EmptyMiddleware()
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		m := c.Request.Method
		t := ginplus.GetTenantID(c)
		u := ginplus.GetUserID(c)
		logrus.Printf("p:%s, m:%s, u:%s, t:%t", p, m, u, t)

		if b, err := enforcer.Enforce(u, t, p, m); err != nil {
			ginplus.ResError(c, errors.WithStack(err))
			return
		} else if !b {
			ginplus.ResError(c, errors.ErrNoPerm)
			return
		}

		c.Next()
	}
}
