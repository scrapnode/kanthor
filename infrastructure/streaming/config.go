package streaming

type ConnectionConfig struct {
	Uri    string       `json:"uri"`
	Stream StreamConfig `json:"stream"`
}

type StreamConfig struct {
	Name     string   `json:"name"`
	Replicas int      `json:"replicas"`
	Subjects []string `json:"subjects"`
	Limits   struct {
		Msgs     int64 `json:"msgs"`
		MsgBytes int32 `json:"msg_bytes"`
		Bytes    int64 `json:"bytes"`
		Age      int64 `json:"age"`
	} `json:"limits"`
}
