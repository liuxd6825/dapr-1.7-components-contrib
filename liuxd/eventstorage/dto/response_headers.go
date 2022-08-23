package dto

type ResponseHeaders struct {
	Values  map[string]string `json:"values"`
	Status  ResponseStatus    `json:"status"`
	Message string            `json:"message"`
	Error   error             `json:"error"`
}

type ResponseStatus int32

const (
	ResponseStatusSuccess        ResponseStatus = iota // 执行成功
	ResponseStatusError                                // 执行错误
	ResponseStatusEventDuplicate                       // 事件已经存在，被重复执行
)

func NewResponseHeaders(status ResponseStatus, err error, values map[string]string) *ResponseHeaders {
	if values == nil {
		values = make(map[string]string)
	}
	if err != nil {
		return NewResponseHeadersError(err, values)
	}
	resp := &ResponseHeaders{
		Status:  status,
		Message: "Success",
		Values:  values,
	}
	return resp
}

func NewResponseHeadersError(err error, values map[string]string) *ResponseHeaders {
	if values == nil {
		values = make(map[string]string)
	}
	resp := &ResponseHeaders{
		Status:  ResponseStatusError,
		Message: err.Error(),
		Values:  values,
	}
	return resp
}

func NewResponseHeadersSuccess(values map[string]string) *ResponseHeaders {
	if values == nil {
		values = make(map[string]string)
	}
	resp := &ResponseHeaders{
		Status:  ResponseStatusSuccess,
		Message: "Success",
		Values:  values,
	}
	return resp
}

func (h *ResponseHeaders) SetHeader(name, value string) *ResponseHeaders {
	h.Values[name] = value
	return h
}

func (h *ResponseHeaders) SetHeaders(data map[string]string) *ResponseHeaders {
	h.Values = data
	return h
}

func (h *ResponseHeaders) SetStatus(data map[string]string) *ResponseHeaders {
	h.Values = data
	return h
}

func (h *ResponseHeaders) GetStatus() ResponseStatus {
	return h.Status
}

func (h *ResponseHeaders) SetError(err error) *ResponseHeaders {
	h.Error = err
	h.Status = ResponseStatusError
	h.Message = err.Error()
	return h
}

func (h *ResponseHeaders) GetError() error {
	return h.Error
}
