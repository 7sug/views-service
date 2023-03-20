package models

type Proxy struct {
	Host string `json:"Host"`
	Port string `json:"Port"`
}

type Response struct {
	CountOfProxy int `json:"CountOfProxy"`
	SuccessCount int `json:"SuccessCount"`
}
