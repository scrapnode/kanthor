package usecase

import "github.com/scrapnode/kanthor/pkg/utils"

func CacheKeyWsAuthenticate(wscId string) string {
	return utils.Key("sdk", "ws", "authenticate", wscId)
}

func CacheKeyApp(wsId, appId string) string {
	return utils.Key("sdk", wsId, appId)
}

func CacheKeyEp(appId, epId string) string {
	return utils.Key("sdk", appId, epId)
}

func CacheKeyEpr(epId, eprId string) string {
	return utils.Key("sdk", epId, eprId)
}
