package container

import (
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	gomock "github.com/golang/mock/gomock"
)

func TestPullImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)

	imgName := "test/img:tag"
	dc.
		EXPECT().
		ImagePull(gomock.Any(), imgName, types.ImagePullOptions{}).
		Return(ioutil.NopCloser(strings.NewReader("")), nil)

	img, err := rt.PullImage(imgName)
	if err != nil {
		t.Errorf("Expected error to be nil but found %v", err)
	}
	expectation := Image{Name: imgName}
	if *img != expectation {
		t.Errorf("Expected %v but found %v", expectation, img)
	}
}

func TestPullImageFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)

	imgName := "test/img:tag"
	dc.
		EXPECT().
		ImagePull(gomock.Any(), imgName, types.ImagePullOptions{}).
		Return(nil, errors.New("an error"))

	_, err := rt.PullImage(imgName)
	if !errors.Is(err, FailedToPullImage) {
		t.Errorf("Expected FailedToPullImage but found %v", err)
	}
}

func TestBuildImage(t *testing.T) {
	// TODO this is pretty rubbish, there is a lot of stuff outside of using the docker
	// API in here.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)

	// Set up temp dir
	workDir, err := ioutil.TempDir("", "test-fold-image-build")
	if err != nil {
		t.Fatalf("failed to create the temporary workspace: %v", err)
	}

	imgName := "test/img:tag"
	dc.
		EXPECT().
		ImageBuild(
			gomock.Any(),
			gomock.Any(),
			types.ImageBuildOptions{Tags: []string{imgName}},
		).
		Return(types.ImageBuildResponse{Body: ioutil.NopCloser(strings.NewReader(""))}, nil)
	rt.BuildImage(&Image{Name: imgName, Src: workDir})
}
