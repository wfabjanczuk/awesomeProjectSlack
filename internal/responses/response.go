package responses

const (
	STATUS_OK    = "OK"
	STATUS_ERROR = "Error"
)

type WSResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
