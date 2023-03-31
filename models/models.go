package models

type Proxy struct {
	Host string `json:"Host"`
	Port string `json:"Port"`
}

type Response struct {
	CountOfProxy int   `json:"CountOfProxy"`
	SuccessCount int64 `json:"SuccessCount"`
}
