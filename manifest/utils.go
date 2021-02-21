package manifest

import (
	"errors"
	"fmt"
	"io"

	"github.com/golang/protobuf/jsonpb"
)

var (
	FailedToWriteJSON = errors.New("failed to convert manifest to JSON")
	FailedToReadJSON  = errors.New("failed to read manifest from JSON")
)

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

func WriteJSON(w io.Writer, m *Manifest) error {
	marshaler := &jsonpb.Marshaler{EmitDefaults: true}
	if err := marshaler.Marshal(w, m); err != nil {
		return FailedToWriteJSON
	}
	return nil
}

func ReadJSON(r io.Reader, m *Manifest) error {
	unmarshaler := &jsonpb.Unmarshaler{}
	if err := unmarshaler.Unmarshal(r, m); err != nil {
		return FailedToReadJSON
	}
	return nil
}
