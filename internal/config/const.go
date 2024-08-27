package config

import (
	"time"

	"github.com/pkg/errors"
)

const (
	CookieSessionName      = "___sid"
	CookieSessionAdminName = "___sid_admin"

	ContextSessionKey          = "context_session_id"
	ContextSessionIsNeedDelete = "context_session_is_need_delete"
	ContextUserUUIDKey         = "context_user_uuid"
	ContextUserLoginKey        = "context_user_login"

	ContextRequestPreviousURIKey = "context_request_previous_uri"
	ContextRequestURIKey         = "context_request_uri"
	ContextRequestBodyKey        = "context_request_body"
	ContextCookiesMapKey         = "context_cookie"

	ContextReturnErrorKey = "context_return_error"

	ContextRequestKey  = "context_request"
	ContextResponseKey = "context_response"

	LayoutDate = "02.01.2006 в 15:04"
)

var (
	DefaultLocation *time.Location
)

func init() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(errors.Wrap(err, "load location"))
	}

	DefaultLocation = loc
}
