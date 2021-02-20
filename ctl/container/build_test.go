package container

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/foldsh/fold/internal/testutils"
	gomock "github.com/golang/mock/gomock"
)

func TestPullImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	fs := &mockFileSystem{}
	rt := mockRuntime(dc, fs)

	imgName := "fold/img:tag"
	dc.
		EXPECT().
		ImagePull(gomock.Any(), imgName, types.ImagePullOptions{}).
		Return(ioutil.NopCloser(strings.NewReader("")), nil)
	dc.
		EXPECT().
		ImageList(gomock.Any(), gomock.Any()).
		Return([]types.ImageSummary{{ID: "1234", RepoTags: []string{imgName}}}, nil)
	dc.
		EXPECT().
		ImageInspectWithRaw(gomock.Any(), "1234").
		Return(
			types.ImageInspect{ID: "1234", Config: &container.Config{WorkingDir: "/fold"}},
			nil,
			nil,
		)

	img, err := rt.PullImage(imgName)
	if err != nil {
		t.Errorf("Expected error to be nil but found %v", err)
	}
	expectation := Image{ID: "1234", Name: imgName, WorkDir: "/fold"}
	if *img != expectation {
		t.Errorf("Expected %v but found %v", expectation, img)
	}
}

func TestPullImageFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	fs := &mockFileSystem{}
	rt := mockRuntime(dc, fs)

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
	fs := &mockFileSystem{}
	rt := mockRuntime(dc, fs)

	// Set up temp dir
	workDir, err := ioutil.TempDir("", "test-fold-image-build")
	if err != nil {
		t.Fatalf("failed to create the temporary workspace: %v", err)
	}

	imgName := "fold/img:tag"
	dc.
		EXPECT().
		ImageBuild(
			gomock.Any(),
			gomock.Any(),
			types.ImageBuildOptions{Tags: []string{imgName}},
		).
		Return(types.ImageBuildResponse{Body: ioutil.NopCloser(strings.NewReader(""))}, nil)
	dc.
		EXPECT().
		ImageList(gomock.Any(), gomock.Any()).
		Return([]types.ImageSummary{{ID: "1234", RepoTags: []string{imgName}}}, nil)
	dc.
		EXPECT().
		ImageInspectWithRaw(gomock.Any(), "1234").
		Return(
			types.ImageInspect{ID: "1234", Config: &container.Config{WorkingDir: "/fold"}},
			nil,
			nil,
		)
	img := &Image{Name: imgName, Src: workDir}
	rt.BuildImage(img)
	if img.ID != "1234" {
		t.Errorf("Expected id to be 1234 but found %s", img.ID)
	}
	if img.WorkDir != "/fold" {
		t.Errorf("Expected id to be /fold but found %s", img.WorkDir)
	}
}

func TestListImages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	fs := &mockFileSystem{}
	rt := mockRuntime(dc, fs)

	dc.
		EXPECT().
		ImageList(
			gomock.Any(),
			gomock.Any(),
		).
		Return([]types.ImageSummary{
			{ID: "0", RepoTags: []string{"foldsh/foo", "bar"}},
			{ID: "1", RepoTags: []string{"foldsh/bar", "fold"}},
			{ID: "2", RepoTags: []string{"foldlocal.abcde.foo", "baz"}},
			{ID: "3", RepoTags: []string{"foo", "bar"}},
		}, nil)

	expectedTags := []string{"foldsh/foo", "foldsh/bar", "foldlocal.abcde.foo"}
	for i, _ := range expectedTags {
		id := fmt.Sprintf("%d", i)
		dc.
			EXPECT().
			ImageInspectWithRaw(gomock.Any(), id).
			Return(
				types.ImageInspect{ID: id, Config: &container.Config{WorkingDir: "/fold"}},
				nil,
				nil,
			)
	}
	imgs, _ := rt.ListImages()
	expectation := []*Image{
		{ID: "0", WorkDir: "/fold", Name: "foldsh/foo"},
		{ID: "1", WorkDir: "/fold", Name: "foldsh/bar"},
		{ID: "2", WorkDir: "/fold", Name: "foldlocal.abcde.foo"},
	}
	testutils.Diff(t, expectation, imgs, "Images should match expectation")
}
