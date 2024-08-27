package utils

import (
	"net/http"
)

func ResponseIsRedirect(status int) bool {
	return status == http.StatusFound ||
		status == http.StatusMovedPermanently ||
		status == http.StatusTemporaryRedirect ||
		status == http.StatusPermanentRedirect
}
