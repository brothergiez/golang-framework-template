package response

type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	Page      int    `json:"page,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Total     int64  `json:"total,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Envelope struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
	Meta    *Meta      `json:"meta,omitempty"`
}

func Success(data any, meta *Meta) Envelope {
	return Envelope{
		Success: true,
		Data:    data,
		Meta:    meta,
	}
}

func Fail(code, message string, details any, meta *Meta) Envelope {
	return Envelope{
		Success: false,
		Error: &ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: meta,
	}
}
