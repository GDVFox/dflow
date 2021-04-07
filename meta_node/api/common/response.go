package common

import (
	"encoding/json"
	"net/http"
)

// Response представляет ответ обработчика.
type Response struct {
	StatusCode int
	Body       []byte
	IsJSON     bool
}

// NewOKResponse возвращает 200 OK, если body ненулевое и 204 No Content, если body нулевое.
func NewOKResponse(body []byte, isJSON bool) *Response {
	if body == nil {
		return &Response{StatusCode: http.StatusNoContent}
	}
	return &Response{
		StatusCode: http.StatusOK,
		Body:       body,
		IsJSON:     isJSON,
	}
}

// NewBadRequestResponse возвращает Response с кодом 400 Bad Request
func NewBadRequestResponse(body []byte) *Response {
	return &Response{
		StatusCode: http.StatusBadRequest,
		Body:       body,
		IsJSON:     true,
	}
}

// NewNotFoundResponse возвращает Response с кодом 404 Not Found
func NewNotFoundResponse(body []byte) *Response {
	return &Response{
		StatusCode: http.StatusNotFound,
		Body:       body,
		IsJSON:     true,
	}
}

// NewConflictResponse возвращает Response с кодом 409 Conflict
func NewConflictResponse(body []byte) *Response {
	return &Response{
		StatusCode: http.StatusConflict,
		Body:       body,
		IsJSON:     true,
	}
}

// NewInternalErrorResponse возвращает Response с кодом 500 Internal Server Error
func NewInternalErrorResponse(body []byte) *Response {
	return &Response{
		StatusCode: http.StatusInternalServerError,
		Body:       body,
		IsJSON:     true,
	}
}

// ErrorCode тип для кодов ошибок
type ErrorCode string

// Возможные коды ошибок.
var (
	BadUnmarshalRequestErrorCode ErrorCode = "bad_unmarshal"
	BadSchemeErrorCode           ErrorCode = "bad_scheme"
	BadActionErrorCode           ErrorCode = "bad_action"
	BadNameErrorCode             ErrorCode = "bad_name"
	NameNotFoundErrorCode        ErrorCode = "name_not_found"
	NameAlreadyExistsErrorCode   ErrorCode = "name_already_exists"
	ETCDErrorCode                ErrorCode = "etcd_error"
)

// ErrorBody тело ответа, при возникновении ошибки.
type ErrorBody struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

// NewErrorBody создает новое тело ответа, содержащего ошибку
func NewErrorBody(c ErrorCode, msg string) []byte {
	errBody := &ErrorBody{
		Code:    c,
		Message: msg,
	}
	data, err := json.Marshal(errBody)
	if err != nil {
		// логируем именно тут для удобства использования в связке с Response.
		// если будем возвращать ошибку, то будет неудобно.
		return nil
	}
	return data
}
