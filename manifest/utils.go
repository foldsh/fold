package manifest

import "fmt"

func HttpMethodFromString(method string) HttpMethod {
	switch method {
	case "GET":
		return HttpMethod_GET
	case "PUT":
		return HttpMethod_PUT
	case "POST":
		return HttpMethod_POST
	case "DELETE":
		return HttpMethod_DELETE
	default:
		// TODO replace with logger.Fatal
		panic(fmt.Sprintf("unsupported http method %s", method))
	}
}
