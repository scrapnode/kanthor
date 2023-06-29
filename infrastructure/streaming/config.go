package streaming

type ConnectionConfig struct {
	Uri    string       `json:"uri" mapstructure:"uri"`
	Stream StreamConfig `json:"stream" mapstructure:"stream"`
}

type StreamConfig struct {
	Name     string   `json:"name" mapstructure:"name"`
	Replicas int      `json:"replicas" mapstructure:"replicas"`
	Subjects []string `json:"subjects" mapstructure:"subjects"`
	Limits   struct {
		Msgs     int64 `json:"msgs" mapstructure:"msgs"`
		MsgBytes int32 `json:"msg_bytes" mapstructure:"msg_bytes"`
		Bytes    int64 `json:"bytes" mapstructure:"bytes"`
		Age      int64 `json:"age" mapstructure:"age"`
	} `json:"limits" mapstructure:"limits"`
}
