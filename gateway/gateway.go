package gateway

type Query struct {
	Search string `json:"_q" form:"_q"`
	Limit  int    `json:"_limit" form:"_limit"`
	Page   int    `json:"_page" form:"_page"`
}