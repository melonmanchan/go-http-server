package statuscodes

import "fmt"

type HTTPStatus struct {
	status int
	msg    string
}

func (h HTTPStatus) ToHeader() string {
	return fmt.Sprintf("HTTP/1.1 %d %s\r\n", h.status, h.msg)
}

var NotFound = HTTPStatus{404, "not found"}
var Ok = HTTPStatus{200, "OK"}
var MethodNotAllowed = HTTPStatus{405, "method not allowed"}
var ServerError = HTTPStatus{500, "internal server error"}
var Teapot = HTTPStatus{418, "I'm a teapot"}
