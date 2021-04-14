package container_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/api/types"
	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/foldsh/fold/ctl/container"
)

func TestContainerStartAndStop(t *testing.T) {
	// TODO improve this by checking the container create config in detail
	rt, dc, _ := setup()
	c := &container.Container{
		Name:  "test",
		Image: container.Image{Name: "fold/test", WorkDir: "/fold"},
		Mounts: []container.Mount{
			{"/home/test/blah/src", "/dst"},
			{"/home/test/blah/foo", "/bar"},
		},
		Environment: map[string]string{"FOLD_SERVICE_NAME": "test"},
	}
	// mfs.On("MkdirAll", "/home/test/blah/src", fs.DIR_PERMISSIONS).Return(nil)
	// mfs.On("MkdirAll", "/home/test/blah/foo", fs.DIR_PERMISSIONS).Return(nil)
	dc.On(
		"ContainerCreate",
		mock.Anything,
		&dockerContainer.Config{
			Image: "fold/test",
			Env: []string{
				"FOLD_STAGE=LOCAL",
				fmt.Sprintf("FOLD_WATCH_DIR=%s", c.Mounts[0].Dst),
				"FOLD_SERVICE_NAME=test",
			},
		},
		mock.Anything,
		mock.Anything,
		mock.Anything,
		"test",
	).Return(
		dockerContainer.ContainerCreateCreatedBody{ID: "testContainerID"},
		nil,
	)
	// dc.On("ContainerStart", mock.Anything, "testContainerID", mock.Anything).Return(nil)

	rt.RunContainer(&container.Network{}, c)

	// assert.Equal(t, "testContainerID", con.ID)

	// dc.On("ContainerStop", mock.Anything, "testContainerID", mock.Anything)
	// rt.StopContainer(con)

	dc.AssertExpectations(t)
	// mfs.AssertExpectations(t)
}

func TestContainerCreateFailure(t *testing.T) {
	rt, dc, _ := setup()
	con := &container.Container{
		Name:   "test",
		Image:  container.Image{Name: "fold/test"},
		Mounts: []container.Mount{{"foo", "bar"}},
	}

	dc.On(
		"ContainerCreate",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(
		dockerContainer.ContainerCreateCreatedBody{},
		errors.New("an error"),
	)

	err := rt.RunContainer(&container.Network{}, con)
	assert.ErrorIs(t, err, container.FailedToCreateContainer)
	dc.AssertExpectations(t)
}

func TestContainerStartAndStopFailure(t *testing.T) {
	rt, dc, _ := setup()
	con := &container.Container{
		Name:   "test",
		Image:  container.Image{Name: "fold/test"},
		Mounts: []container.Mount{{"foo", "bar"}},
	}

	dc.On(
		"ContainerCreate",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(
		dockerContainer.ContainerCreateCreatedBody{},
		nil,
	)

	dc.On(
		"ContainerStart",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(
		errors.New("an error"),
	)

	err := rt.RunContainer(&container.Network{}, con)
	assert.ErrorIs(t, err, container.FailedToStartContainer)
	dc.On(
		"ContainerStop",
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(
		errors.New("an error"),
	)

	err = rt.StopContainer(con)
	assert.ErrorIs(t, err, container.FailedToStopContainer)
	dc.AssertExpectations(t)
}

func TestContainerRemove(t *testing.T) {
	rt, dc, _ := setup()
	con := &container.Container{
		ID:     "testContainerID",
		Name:   "test",
		Image:  container.Image{Name: "fold/test"},
		Mounts: []container.Mount{{"foo", "bar"}},
	}

	dc.
		On("ContainerRemove", mock.Anything, "testContainerID", mock.Anything).Return(nil)
	err := rt.RemoveContainer(con)
	assert.Nil(t, err)
	dc.AssertExpectations(t)
}

func TestContainerRemoveFailure(t *testing.T) {
	rt, dc, _ := setup()
	con := &container.Container{
		ID:     "testContainerID",
		Name:   "test",
		Image:  container.Image{Name: "fold/test"},
		Mounts: []container.Mount{{"foo", "bar"}},
	}

	dc.
		On("ContainerRemove", mock.Anything, "testContainerID", mock.Anything).
		Return(errors.New("an error"))
	err := rt.RemoveContainer(con)
	assert.ErrorIs(t, err, container.FailedToRemoveContainer)
	dc.AssertExpectations(t)
}

func TestContainerJoinAndLeaveNetwork(t *testing.T) {
	rt, dc, _ := setup()
	net := &container.Network{Name: "testNet", ID: "testNetID"}
	con := &container.Container{
		Name:   "testCon",
		ID:     "testConID",
		Image:  container.Image{Name: "fold/test"},
		Mounts: []container.Mount{{"foo", "bar"}},
	}

	// Happy
	dc.On("NetworkConnect", mock.Anything, "testNetID", "testConID", mock.Anything).Return(nil)
	rt.AddToNetwork(net, con)

	dc.On("NetworkDisconnect", mock.Anything, "testNetID", "testConID", false).Return(nil)
	rt.RemoveFromNetwork(net, con)
	dc.AssertExpectations(t)
}

func TestContainerJoinAndLeaveNetworkFailure(t *testing.T) {
	rt, dc, _ := setup()
	net := &container.Network{Name: "testNet", ID: "testNetID"}
	con := &container.Container{
		Name:   "testCon",
		ID:     "testConID",
		Image:  container.Image{Name: "fold/test"},
		Mounts: []container.Mount{{"foo", "bar"}},
	}

	dc.
		On("NetworkConnect", mock.Anything, "testNetID", "testConID", mock.Anything).
		Return(errors.New("an error"))
	err := rt.AddToNetwork(net, con)
	assert.ErrorIs(t, err, container.FailedToJoinNetwork)

	dc.
		On("NetworkDisconnect", mock.Anything, "testNetID", "testConID", false).
		Return(errors.New("an error"))
	err = rt.RemoveFromNetwork(net, con)
	assert.ErrorIs(t, err, container.FailedToLeaveNetwork)
	dc.AssertExpectations(t)
}

func TestListAllFoldContainers(t *testing.T) {
	rt, dc, _ := setup()
	dc.
		On("ContainerList", mock.Anything, mock.Anything).
		Return([]types.Container{
			containerResponse("a", "/foo", "/bar"),
			containerResponse("b", "/fold.foo", "/bar"),
			containerResponse("c", "/foo", "/fold.bar"),
			containerResponse("d", "/fold.baz"),
			containerResponse("e", "/bar"),
		}, nil)

	cs, err := rt.AllContainers()
	assert.Nil(t, err)
	expectation := []*container.Container{
		{
			ID:     "b",
			Name:   "fold.foo",
			Image:  container.Image{Name: "test"},
			Mounts: []container.Mount{container.Mount{"src", "dst"}},
		},
		{
			ID:     "c",
			Name:   "fold.bar",
			Image:  container.Image{Name: "test"},
			Mounts: []container.Mount{container.Mount{"src", "dst"}},
		},
		{
			ID:     "d",
			Name:   "fold.baz",
			Image:  container.Image{Name: "test"},
			Mounts: []container.Mount{container.Mount{"src", "dst"}},
		},
	}
	diffContainers(t, expectation, cs)
	dc.AssertExpectations(t)
}

func TestListAllFoldContainersFailure(t *testing.T) {
	rt, dc, _ := setup()
	dc.
		On("ContainerList", mock.Anything, mock.Anything).
		Return([]types.Container{}, errors.New("an error"))

	_, err := rt.AllContainers()
	assert.ErrorIs(t, err, container.DockerEngineError)
	dc.AssertExpectations(t)
}

func TestGetContainer(t *testing.T) {
	rt, dc, _ := setup()
	dc.
		On("ContainerList", mock.Anything, mock.Anything).
		Return([]types.Container{
			containerResponse("a", "/foo", "/bar"),
			containerResponse("b", "/fold.foo", "/bar"),
			containerResponse("c", "/foo", "/fold.bar"),
			containerResponse("d", "/fold.baz"),
			containerResponse("e", "/bar"),
		}, nil)

	c, err := rt.GetContainer("foo")
	assert.Nil(t, err)
	expectation := &container.Container{
		ID: "b", Name: "fold.foo", Image: container.Image{Name: "test"}, Mounts: []container.Mount{{"src", "dst"}},
	}
	diffContainers(t, expectation, c)
	dc.AssertExpectations(t)
}

func TestContainerLogs(t *testing.T) {
	rt, dc, _ := setup()
	con := &container.Container{
		Name:   "testCon",
		ID:     "testConID",
		Image:  container.Image{Name: "fold/test"},
		Mounts: []container.Mount{{"foo", "bar"}},
	}

	dc.
		On("ContainerLogs", mock.Anything, con.ID, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Follow:     true,
		}).
		Return(
			ioutil.NopCloser(
				bytes.NewReader(
					[]byte{
						// The first byte of the header says which message stream this is from
						0x01,
						// The next 3 bytes are always 0
						0x00, 0x00, 0x00,
						// The next 4 bytes are the message size in big endian layout (7 in this case)
						0x00, 0x00, 0x00, 0x07,
						// Then we have the message body, in this case the text 'foo bar'
						0x66, 0x6f, 0x6f, 0x20, 0x62, 0x61, 0x72,
					},
				),
			), nil,
		)
	ls, _ := rt.ContainerLogs(con)
	var buf bytes.Buffer
	err := ls.Stream(&buf)
	require.Nil(t, err)
	ls.Stop()
	assert.Equal(t, "foo bar", buf.String())
	dc.AssertExpectations(t)
}

func containerResponse(id string, names ...string) types.Container {
	return types.Container{
		ID:     id,
		Names:  names,
		Image:  "test",
		Mounts: []types.MountPoint{{Source: "src", Destination: "dst"}},
	}
}

func diffContainers(t *testing.T, expectation, actual interface{}) {
	if diff := cmp.Diff(
		expectation,
		actual,
		cmpopts.IgnoreUnexported(container.Container{}),
	); diff != "" {
		t.Errorf("Expected container lists to be equal(-want +got):\n%s", diff)
	}
}
