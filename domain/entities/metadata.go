package entities

const (
	MetaMsgId = "kanthor.msg.id"
	MetaEpId  = "kanthor.ep.id"
	MetaEprId = "kanthor.epr.id"
	MetaReqId = "kanthor.req.id"
	MetaResId = "kanthor.res.id"
	MetaAttId = "kanthor.att.id"
)

type Metadata map[string]string

func (meta Metadata) Get(key string) string {
	return meta[key]
}

func (meta Metadata) Set(key, value string) {
	meta[key] = value
}

func (meta Metadata) Merge(target map[string]string) {
	for key, value := range target {
		meta[key] = value
	}
}
