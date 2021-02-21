package manifest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/foldsh/fold/internal/testutils"
	"github.com/foldsh/fold/manifest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
			{HttpMethod: manifest.HttpMethod_GET, Handler: "get", PathSpec: "/get/:var"},
			{HttpMethod: manifest.HttpMethod_PUT, Handler: "put", PathSpec: "/put/:var"},
			{HttpMethod: manifest.HttpMethod_POST, Handler: "post", PathSpec: "/post/:var"},
			{HttpMethod: manifest.HttpMethod_DELETE, Handler: "delete", PathSpec: "/delete/:var"},
			{HttpMethod: manifest.HttpMethod_PATCH, Handler: "patch", PathSpec: "/patch/:var"},
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
			map[string]interface{}{"httpMethod": "GET", "handler": "get", "pathSpec": "/get/:var"},
			map[string]interface{}{"httpMethod": "PUT", "handler": "put", "pathSpec": "/put/:var"},
			map[string]interface{}{
				"httpMethod": "POST",
				"handler":    "post",
				"pathSpec":   "/post/:var",
			},
			map[string]interface{}{
				"httpMethod": "DELETE",
				"handler":    "delete",
				"pathSpec":   "/delete/:var",
			},
			map[string]interface{}{
				"httpMethod": "PATCH",
				"handler":    "patch",
				"pathSpec":   "/patch/:var",
			},
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
	if diff := cmp.Diff(
		m,
		result,
		cmpopts.IgnoreUnexported(manifest.Manifest{}, manifest.Version{}, manifest.BuildInfo{}, manifest.Route{}),
	); diff != "" {
		t.Errorf("Manifest read from JSON does not match exepctation(-want +got):\n%s", diff)
	}
}
