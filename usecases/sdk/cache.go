package sdk

import "github.com/scrapnode/kanthor/pkg/utils"

func CacheKeyWsAuthenticate(wscId string) string {
	return utils.Key("sdk", "workspace", "authenticate", wscId)
}

func CacheKeyApp(wsId, appId string) string {
	return utils.Key("sdk", wsId, appId)
}

func CacheKeyEp(wsId, appId, epId string) string {
	return utils.Key("sdk", wsId, appId, epId)
}

func CacheKeyEpr(wsId, appId, epId, eprId string) string {
	return utils.Key("sdk", wsId, appId, epId, eprId)
}
