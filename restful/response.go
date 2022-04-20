package restful

type Response struct {
	Status  *int    `json:"status,omitempty"`
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
	Tips    *string `json:"tips,omitempty"`
	Data    any     `json:"data,omitempty"`
}
