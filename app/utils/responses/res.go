package responses

import (
	"github.com/grushka14/S4/domain/token"
)

// Response type
type Response struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// UsersResponse type
type UsersResponse struct {
	Code   int      `json:"code"`
	Status string   `json:"status"`
	Users  []string `json:"users"`
}

// GetFilesResponse type
type GetFilesResponse struct {
	Code   int                 `json:"code"`
	Status string              `json:"status"`
	Files  []map[string]string `json:"files"`
}

// ReadFilesResponse type
type ReadFilesResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	File   []byte `json:"file"`
}

// LoginResponse type
type LoginResponse struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Token  token.Token `json:"token"`
}

// New400 func
func New400(m string) *Response {
	return &Response{
		Code:    400,
		Status:  "bad_request",
		Message: m,
	}
}

// New200 func
func New200(m string) *Response {
	return &Response{
		Code:    200,
		Status:  "ok",
		Message: m,
	}
}

// New500 func
func New500(m string) *Response {
	return &Response{
		Code:    500,
		Status:  "server_error",
		Message: m,
	}
}

// New401 func
func New401(m string) *Response {
	return &Response{
		Code:    401,
		Status:  "unauthorized",
		Message: m,
	}
}

// New204 func
func New204(m string) *Response {
	return &Response{
		Code:    204,
		Status:  "resource deleted",
		Message: m,
	}
}

// New201 func
func New201(m string) *Response {
	return &Response{
		Code:    204,
		Status:  "created",
		Message: m,
	}
}

// New200GetUsers func
func New200GetUsers(u []string) *UsersResponse {
	return &UsersResponse{
		Code:   200,
		Status: "OK",
		Users:  u,
	}
}

// New200GetFiles func
func New200GetFiles(u []map[string]string) *GetFilesResponse {
	return &GetFilesResponse{
		Code:   200,
		Status: "OK",
		Files:  u,
	}
}

// New200GetFiles func
func New200ReadFile(f []byte) *ReadFilesResponse {
	return &ReadFilesResponse{
		Code:   200,
		Status: "OK",
		File:   f,
	}
}

// New200Login func
func New200Login(t token.Token) *LoginResponse {
	return &LoginResponse{
		Code:   200,
		Status: "OK",
		Token:  t,
	}
}
