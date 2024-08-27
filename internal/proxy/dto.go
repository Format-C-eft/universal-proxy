package proxy

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type Request struct {
	Method     string
	RequestURI string
	Headers    http.Header
	Cookies    []*http.Cookie
	Body       []byte
	User       User
}

type User struct {
	SessionID string
	Login     string
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Cookies    []http.Cookie
	Body       []byte
}

type proxy struct {
	IP         string
	ExternalIP string
	Client     *resty.Client
	SessionID  string
	Login      string
	LastURL    string
	LastUsed   time.Time
	IsActive   bool
	CountErr   int
	LastErr    error
}

type Info struct {
	Address         string `json:"address"`
	AddressExternal string `json:"address_external"`
	SessionID       string `json:"session_id"`
	Login           string `json:"login"`
	LastURL         string `json:"last_url"`
	LastUsed        string `json:"last_used"`
	IsActive        bool   `json:"is_active"`
	IsUsedNow       bool   `json:"is_used_now"`
	LastErr         string `json:"last_err"`
	CountErr        int    `json:"count_err"`
}
