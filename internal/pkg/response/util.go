package response

type Error struct {
	httpCode  int
	errorCode int
	errorMsg  string
}

type Response struct {
	ErrorCode int         `json:"error_code"`
	Data      interface{} `json:"data"`
	Message   string      `json:"message"`
}
