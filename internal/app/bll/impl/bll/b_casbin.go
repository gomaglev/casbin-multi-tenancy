package bll

import (
	"context"

	"gin-casbin/internal/app/config"
	"gin-casbin/pkg/logger"

	"github.com/casbin/casbin/v2"
)

var chCasbinPolicy chan *chCasbinPolicyItem

type chCasbinPolicyItem struct {
	ctx context.Context
	e   *casbin.SyncedEnforcer
}

func init() {
	chCasbinPolicy = make(chan *chCasbinPolicyItem, 1)
	go func() {
		for item := range chCasbinPolicy {
			err := item.e.LoadPolicy()
			if err != nil {
				logger.Errorf(item.ctx, "The load casbin policy error: %s", err.Error())
			}
		}
	}()
}

// LoadCasbinPolicy
func LoadCasbinPolicy(ctx context.Context, e *casbin.SyncedEnforcer) {
	if !config.C.Casbin.Enable {
		return
	}

	if len(chCasbinPolicy) > 0 {
		logger.Infof(ctx, "The load casbin policy is already in the wait queue")
		return
	}

	chCasbinPolicy <- &chCasbinPolicyItem{
		ctx: ctx,
		e:   e,
	}
}
