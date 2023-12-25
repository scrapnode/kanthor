package gateway

type Query struct {
	Search string   `json:"_q" form:"_q"`
	Limit  int      `json:"_limit" form:"_limit"`
	Page   int      `json:"_page" form:"_page"`
	Id     []string `json:"id" form:"id"`
	Start  int64    `json:"_start" form:"_start"`
	End    int64    `json:"_end" form:"_end"`
}
