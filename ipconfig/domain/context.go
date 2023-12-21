package domain

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type IpconfigContext struct {
	Ctx          *context.Context
	AppCtx       *app.RequestContext
	ClientConext *ClientConext
}
type ClientConext struct {
	IP string `json:"ip"`
}

func BuildIpConfContext(c *context.Context, ctx *app.RequestContext) *IpconfigContext {
	ipConfConext := &IpconfigContext{
		Ctx:          c,
		AppCtx:       ctx,
		ClientConext: &ClientConext{},
	}
	return ipConfConext
}
