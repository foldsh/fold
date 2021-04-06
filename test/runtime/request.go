package runtime

import "net/http"

type ResponseWriter struct {
	statusCode int
	headers    http.Header
	body       []byte
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{headers: make(map[string][]string)}
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.headers
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.body = b
	return len(b), nil
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) Status() int {
	return rw.statusCode
}

func (rw *ResponseWriter) String() string {
	return string(rw.body)
}
