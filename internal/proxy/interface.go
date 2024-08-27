package proxy

type Store interface {
	ExecuteRequest(request *Request) (*Response, error)
	GetProxyInfo() []Info
}
