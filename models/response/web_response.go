package response

type WebResponse struct {
	Code     int         `json:"code"`
	Status   string      `json:"status"`
	Messages *[]string   `json:"message"`
	Content  interface{} `json:"content"`
}

func NewWebResponse(data interface{}, messages ...string) *WebResponse {
	return &WebResponse{
		Code:     200,
		Status:   "Ok",
		Messages: &messages,
		Content:  data,
	}
}

func NewErrorWebResponse(code int, status string, messages ...string) *WebResponse {
	return &WebResponse{
		Code:     code,
		Status:   status,
		Messages: &messages,
		Content:  nil,
	}
}
