package entities

var (
	MetaEpId  = "kanthor.ep.id"
	MetaEprId = "kanthor.epr.id"
)

type Metadata map[string]string

func (meta Metadata) Get(key string) string {
	return meta[key]
}

func (meta Metadata) Set(key, value string) {
	meta[key] = value
}

func (meta Metadata) Merge(src map[string]string) {
	for key, value := range src {
		meta[key] = value
	}
}
