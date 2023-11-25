package ds

import "fmt"

func ReqKey(msgId, epId string) string {
	return fmt.Sprintf("%s/%s/req", msgId, epId)
}

func ResKey(msgId, epId string) string {
	return fmt.Sprintf("%s/%s/res", msgId, epId)
}
