package model

import "github.com/google/wire"

// ModelSet
var ModelSet = wire.NewSet(
	RoleSet,
	TransSet,
	UserRoleSet,
	UserSet,
	TenantSet,
	UserTenantSet,
	TenantAdministratorSet,
)
