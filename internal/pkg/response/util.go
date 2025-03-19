package response

type Response struct {
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
}
