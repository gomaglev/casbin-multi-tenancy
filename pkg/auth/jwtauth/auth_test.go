package jwtauth

import (
	"context"
	"testing"

	"gin-casbin/pkg/auth/jwtauth/store/buntdb"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	store, err := buntdb.NewStore(":memory:")
	assert.Nil(t, err)

	jwtAuth := New(store)

	defer jwtAuth.Release()

	ctx := context.Background()
	userID := "test"
	tenantID := "tenant"
	token, err := jwtAuth.GenerateToken(ctx, userID, tenantID)
	assert.Nil(t, err)
	assert.NotNil(t, token)

	id, tid, err := jwtAuth.ParseUserID(ctx, token.GetAccessToken())
	assert.Nil(t, err)
	assert.Equal(t, userID, id)
	assert.Equal(t, tenantID, tid)

	err = jwtAuth.DestroyToken(ctx, token.GetAccessToken())
	assert.Nil(t, err)

	id, tid, err = jwtAuth.ParseUserID(ctx, token.GetAccessToken())
	assert.NotNil(t, err)
	assert.EqualError(t, err, "invalid token")
	assert.Empty(t, id)
	assert.Empty(t, tid)
}
