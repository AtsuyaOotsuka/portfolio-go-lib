package atylabapi

import (
	"context"
	"time"
)

type ApiCtxSvc struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewApiCtxSvc() *ApiCtxSvc {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return &ApiCtxSvc{
		Ctx:    ctx,
		Cancel: cancel,
	}
}

func NewApiCtxSvcWithTimeout(timeout time.Duration) *ApiCtxSvc {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &ApiCtxSvc{
		Ctx:    ctx,
		Cancel: cancel,
	}
}
