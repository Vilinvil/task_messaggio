package models

//easyjson:json
type MessageStatistic struct {
	Total   int `json:"total"`
	Handled int `json:"handled"`
}

type ResponseMessageStatistic struct {
	*MessageStatistic
	status int
}

func NewResponseMessageStatistic(status int, messageStatistic *MessageStatistic) *ResponseMessageStatistic {
	return &ResponseMessageStatistic{
		MessageStatistic: messageStatistic,
		status:           status,
	}
}

func (r *ResponseMessageStatistic) Status() int {
	return r.status
}
