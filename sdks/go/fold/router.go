package fold

type Router struct{}

type Request struct {
	HttpMethod  string
	Handler     string
	Path        string
	Body        map[string]interface{}
	Headers     map[string][]string
	PathParams  map[string]string
	QueryParams map[string][]string
}

type Response struct {
	StatusCode int
	Body       map[string]interface{}
	Headers    map[string][]string
}

type Handler func(*Request, *Response)
