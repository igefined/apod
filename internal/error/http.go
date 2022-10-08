package error

import (
	"encoding/json"
	"errors"
)

var (
	Unauthorized          = errors.New("unauthorized")
	BadRequest            = errors.New("bad request")
	NotFound              = errors.New("not found")
	InternalSystem        = errors.New("internal system error")
	ItemIsExists          = errors.New("item is exists")
	InvalidRefreshToken   = errors.New("invalid refresh token")
	InvalidAccessToken    = errors.New("invalid access token")
	UserIsExistsWithEmail = errors.New("user is exists with email")
	ImageNotFound         = errors.New("image is not exists")
)

type HttpRestError interface {
	Status() int
	Error() string
	Causes() interface{}
	Marshal() []byte
}

type HttpError struct {
	ErrStatus int         `json:"status"`
	ErrError  string      `json:"error"`
	ErrCauses interface{} `json:"-"`
}

func (h *HttpError) Status() int {
	return h.ErrStatus
}

func (h *HttpError) Error() string {
	return h.ErrError
}

func (h *HttpError) Causes() interface{} {
	return h.ErrCauses
}

func (h *HttpError) Marshal() []byte {
	marshal, err := json.Marshal(h)
	if err != nil {
		return nil
	}

	return marshal
}
