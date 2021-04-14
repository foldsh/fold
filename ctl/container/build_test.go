package container_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/foldsh/fold/ctl/container"
)

func TestPullImage(t *testing.T) {
	rt, dc, fs := setup()

	imgName := "fold/img:tag"

	dc.On(
		"ImagePull",
		mock.Anything,
		imgName,
		types.ImagePullOptions{},
	).Return(
		ioutil.NopCloser(strings.NewReader("")),
		nil,
	)
	dc.On(
		"ImageList",
		mock.Anything,
		mock.Anything,
	).Return(
		[]types.ImageSummary{{ID: "1234", RepoTags: []string{imgName}}},
		nil,
	)
	dc.On("ImageInspectWithRaw", mock.Anything, "1234").Return(
		types.ImageInspect{ID: "1234", Config: &dockerContainer.Config{WorkingDir: "/fold"}},
		nil,
		nil,
	)

	img, err := rt.PullImage(imgName)
	dc.AssertExpectations(t)
	fs.AssertExpectations(t)
	assert.Nil(t, err)
	expectation := container.Image{ID: "1234", Name: imgName, WorkDir: "/fold"}
	assert.Equal(t, expectation, *img)
}

func TestPullImageFailure(t *testing.T) {
	rt, dc, _ := setup()

	imgName := "test/img:tag"
	dc.On(
		"ImagePull",
		mock.Anything,
		imgName,
		types.ImagePullOptions{},
	).Return(
		nil,
		errors.New("an error"),
	)

	_, err := rt.PullImage(imgName)
	assert.ErrorIs(t, err, container.FailedToPullImage)
	dc.AssertExpectations(t)
}

func TestBuildImage(t *testing.T) {
	// TODO this is pretty rubbish, there is a lot of stuff outside of using the docker
	// API in here.
	rt, dc, _ := setup()

	// Set up temp dir
	workDir, err := ioutil.TempDir("", "test-fold-image-build")
	if err != nil {
		t.Fatalf("failed to create the temporary workspace: %v", err)
	}

	imgName := "fold/img:tag"
	dc.On(
		"ImageBuild",
		mock.Anything,
		mock.Anything,
		types.ImageBuildOptions{Tags: []string{imgName}},
	).Return(
		types.ImageBuildResponse{Body: ioutil.NopCloser(strings.NewReader(""))},
		nil,
	)
	dc.On(
		"ImageList",
		mock.Anything,
		mock.Anything,
	).Return(
		[]types.ImageSummary{{ID: "1234", RepoTags: []string{imgName}}},
		nil,
	)
	dc.On(
		"ImageInspectWithRaw",
		mock.Anything,
		"1234",
	).Return(
		types.ImageInspect{ID: "1234", Config: &dockerContainer.Config{WorkingDir: "/fold"}},
		nil,
		nil,
	)
	img := &container.Image{Name: imgName, Src: workDir}
	rt.BuildImage(img)
	assert.Equal(t, "1234", img.ID)
	assert.Equal(t, "/fold", img.WorkDir)
	dc.AssertExpectations(t)
}

func TestListImages(t *testing.T) {
	rt, dc, _ := setup()

	dc.
		On("ImageList",
			mock.Anything,
			mock.Anything,
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
			On("ImageInspectWithRaw", mock.Anything, id).
			Return(
				types.ImageInspect{ID: id, Config: &dockerContainer.Config{WorkingDir: "/fold"}},
				nil,
				nil,
			)
	}
	imgs, _ := rt.ListImages()
	expectation := []*container.Image{
		{ID: "0", WorkDir: "/fold", Name: "foldsh/foo"},
		{ID: "1", WorkDir: "/fold", Name: "foldsh/bar"},
		{ID: "2", WorkDir: "/fold", Name: "foldlocal.abcde.foo"},
	}
	assert.Equal(t, expectation, imgs)
	dc.AssertExpectations(t)
}
