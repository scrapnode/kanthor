package msgbroker

type Config struct {
	Uri      string         `json:"uri"`
	Stream   ConfigStream   `json:"stream"`
	Consumer ConfigConsumer `json:"consumer"`
}

type ConfigStream struct {
	Name     string `json:"name"`
	Replicas int    `json:"replicas"`
	Subject  string `json:"subject"`
	Limits   struct {
		Msgs     int64 `json:"msgs"`
		MsgBytes int32 `json:"msg_bytes"`
		Bytes    int64 `json:"bytes"`
		Age      int64 `json:"age"`
	} `json:"limits"`
}

type ConfigConsumer struct {
	Name          string `json:"name"`
	FilterSubject string `json:"filter_subject"`
	Temporary     bool   `json:"temporary"`
	MaxRetry      int    `json:"max_retry"`
}
