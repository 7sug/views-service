package config

type Settings struct {
	Port                 int    `json:"port"`
	AcceptHeader         string `json:"accept_header"`
	AcceptEncodingHeader string `json:"accept_encoding_header"`
	XRequestedWithHeader string `json:"x_requested_with_header"`
}
