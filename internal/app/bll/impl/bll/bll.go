package bll

import "github.com/google/wire"

// BllSet bll注入
var BllSet = wire.NewSet(
	LoginSet,
	RoleSet,
	UserSet,
	TenantSet,
)
