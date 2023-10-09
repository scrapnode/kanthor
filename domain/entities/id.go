package entities

import (
	"fmt"

	"github.com/scrapnode/kanthor/pkg/suid"
)

var (
	IdNsWs  = "ws"
	IdNsWst = "wst"
	IdNsWsc = "wsc"
	IdNsApp = "app"
	IdNsEp  = "ep"
	IdNsEpr = "epr"
	IdNsMsg = "msg"
	IdNsReq = "req"
	IdNsRes = "res"
)

func Id(ns, id string) string {
	return fmt.Sprintf("%s_%s", ns, id)
}

func WsId() string {
	return suid.New(IdNsWs)
}

func WstId() string {
	return suid.New(IdNsWst)
}

func WscId() string {
	return suid.New(IdNsWsc)
}

func AppId() string {
	return suid.New(IdNsApp)
}

func EpId() string {
	return suid.New(IdNsEp)
}

func EprId() string {
	return suid.New(IdNsEpr)
}

func MsgId() string {
	return suid.New(IdNsMsg)
}

func ReqId() string {
	return suid.New(IdNsReq)
}

func ResId() string {
	return suid.New(IdNsRes)
}
