package injector

import (
	"gin-casbin/internal/app/config"
	"gin-casbin/internal/app/middleware"
	"gin-casbin/internal/app/router"

	"github.com/gin-gonic/gin"
	i18n "github.com/suisrc/gin-i18n"
	"golang.org/x/text/language"
)

// InitGinEngine
func InitGinEngine(r router.IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)

	app := gin.New()
	app.NoMethod(middleware.NoMethodHandler())
	app.NoRoute(middleware.NoRouteHandler())

	bundle := i18n.NewBundle(
		language.Chinese,
		"configs/i18n/active.zh-CN.toml",
		"configs/i18n/active.en-US.toml",
		"configs/i18n/active.ja-JP.toml",
	)

	// I18n
	app.Use(i18n.Serve(bundle))

	r.Register(app)

	return app
}
