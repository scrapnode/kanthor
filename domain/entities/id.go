package entities

import "github.com/scrapnode/kanthor/pkg/utils"

func WsId() string {
	return utils.ID("ws")
}

func WstId() string {
	return utils.ID("wst")
}

func WscId() string {
	return utils.ID("wsc")
}

func AppId() string {
	return utils.ID("app")
}

func EpId() string {
	return utils.ID("ep")
}

func EprId() string {
	return utils.ID("epr")
}

func MsgId() string {
	return utils.ID("msg")
}

func ReqId() string {
	return utils.ID("req")
}

func ResId() string {
	return utils.ID("res")
}

func AttId() string {
	return utils.ID("att")
}
