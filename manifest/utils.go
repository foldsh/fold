package manifest

import "fmt"

type InvalidHTTPMethod struct {
	HTTPMethod string
}

func (ihm InvalidHTTPMethod) Error() string {
	return fmt.Sprintf("invalid HTTP method %s", ihm.HTTPMethod)
}

func HttpMethodFromString(method string) (HttpMethod, error) {
	if value, ok := HttpMethod_value[method]; ok {
		return HttpMethod(value), nil
	} else {
		return -1, InvalidHTTPMethod{method}
	}
}
