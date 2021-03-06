// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package injector

import (

	// "gin-casbin/internal/app/api/mock"

	"gin-casbin/internal/app/module/adapter"

	"github.com/google/wire"

	gormModel "gin-casbin/internal/app/model/impl/gorm/model"
)

// BuildInjector
func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		InitGormDB,
		gormModel.ModelSet,
		InitAuth,
		InitCasbin,
		InitGinEngine,
		adapter.CasbinAdapterSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
