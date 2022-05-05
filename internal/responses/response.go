package responses

const (
	WS_STATUS_OK = "OK"
)

type WSResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
