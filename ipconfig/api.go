package ipconfig

import (
	"context"

	"hello/plato_mjr/ipconfig/domain"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type IpconfigResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

func GetIpList(c context.Context, ctx *app.RequestContext) {
	defer func() {
		if err := recover(); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"err": err})
		}
	}()
	// Step0: 构建客户请求信息
	ipConfCtx := domain.BuildIpConfContext(&c, ctx)
	// Step1: 进行ip调度
	eds := domain.Dispatch(ipConfCtx)
	// Step2: 根据得分取top5返回
	ipConfCtx.AppCtx.JSON(consts.StatusOK, packRes(eds))
}
