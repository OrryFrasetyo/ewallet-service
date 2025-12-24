package handler

type WebResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // omitempty: kalau null, gak usah tampil
	Error   string      `json:"error,omitempty"`
}
