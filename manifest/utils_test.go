package manifest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/foldsh/fold/internal/testutils"
	"github.com/foldsh/fold/manifest"
)

var (
	m = &manifest.Manifest{
		Name:    "test",
		Version: &manifest.Version{Major: 1, Minor: 2, Patch: 3},
		BuildInfo: &manifest.BuildInfo{
			Maintainer: "tom",
			Image:      "image",
			Tag:        "tag",
			Path:       "./build/path",
		},
		Routes: []*manifest.Route{
			{HttpMethod: manifest.FoldHTTPMethod_GET, Route: "/get/:var"},
			{HttpMethod: manifest.FoldHTTPMethod_PUT, Route: "/put/:var"},
			{HttpMethod: manifest.FoldHTTPMethod_POST, Route: "/post/:var"},
			{HttpMethod: manifest.FoldHTTPMethod_DELETE, Route: "/delete/:var"},
			{HttpMethod: manifest.FoldHTTPMethod_PATCH, Route: "/patch/:var"},
		},
	}
	j = map[string]interface{}{
		"name": "test",
		"version": map[string]interface{}{
			"major": float64(1),
			"minor": float64(2),
			"patch": float64(3),
		},
		"buildInfo": map[string]interface{}{
			"maintainer": "tom",
			"image":      "image",
			"tag":        "tag",
			"path":       "./build/path",
		},
		"routes": []interface{}{
			map[string]interface{}{"httpMethod": "GET", "route": "/get/:var"},
			map[string]interface{}{"httpMethod": "PUT", "route": "/put/:var"},
			map[string]interface{}{"httpMethod": "POST", "route": "/post/:var"},
			map[string]interface{}{"httpMethod": "DELETE", "route": "/delete/:var"},
			map[string]interface{}{"httpMethod": "PATCH", "route": "/patch/:var"},
		},
	}
)

func TestWriteJSON(t *testing.T) {
	buf := &bytes.Buffer{}
	manifest.WriteJSON(buf, m)

	var actual map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &actual); err != nil {
		t.Fatalf("%+v", err)
	}
	fmt.Printf("%+v", actual)

	testutils.Diff(t, j, actual, "JSON did not match expectation")
}

func TestReadJSON(t *testing.T) {
	bs, err := json.Marshal(j)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	buf := bytes.NewBuffer(bs)

	result := &manifest.Manifest{}
	if err := manifest.ReadJSON(buf, result); err != nil {
		t.Fatalf("%+v", err)
	}
	testutils.DiffManifest(t, m, result)
}
