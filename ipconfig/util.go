package ipconfig

import "hello/plato_mjr/ipconfig/domain"

func packRes(eps []*domain.Endpoint) *IpconfigResponse {
	return &IpconfigResponse{
		Message: "success",
		Code:    0,
		Data:    eps,
	}
}
