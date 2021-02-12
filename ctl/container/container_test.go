package container

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/foldsh/fold/logging"
)

var any = gomock.Any()

func TestContainerStartAndStop(t *testing.T) {
	// TODO improve this by checking the container create config in detail
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)
	con := &Container{
		Name:   "test",
		Image:  Image{Name: "fold/test", WorkDir: "/fold"},
		Mounts: []Mount{{"foo", "bar"}},
	}
	dc.
		EXPECT().
		ContainerCreate(
			any, &container.Config{Image: "fold/test", Env: []string{
				"FOLD_STAGE=LOCAL",
				fmt.Sprintf("FOLD_WATCH_DIR=%s", con.Mounts[0].Dst),
			}}, any, any, any, "test",
		).
		Return(container.ContainerCreateCreatedBody{ID: "testContainerID"}, nil)
	dc.
		EXPECT().
		ContainerStart(any, "testContainerID", any)
	rt.RunContainer(&Network{}, con)
	if con.ID != "testContainerID" {
		t.Errorf("Expected container ID to be set after start")
	}
	dc.
		EXPECT().
		ContainerStop(any, "testContainerID", any)
	rt.StopContainer(con)
}

func TestContainerStartAndStopFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)
	con := &Container{
		Name:   "test",
		Image:  Image{Name: "fold/test"},
		Mounts: []Mount{{"foo", "bar"}},
	}

	dc.
		EXPECT().
		ContainerCreate(any, any, any, any, any, any).
		Return(container.ContainerCreateCreatedBody{}, errors.New("an error"))
	err := rt.RunContainer(&Network{}, con)
	if !errors.Is(err, FailedToCreateContainer) {
		t.Errorf("Expected FailedToCreateContainer but found %v", err)
	}
	dc.
		EXPECT().
		ContainerCreate(any, any, any, any, any, any).
		Return(container.ContainerCreateCreatedBody{}, nil)
	dc.
		EXPECT().
		ContainerStart(any, any, any).
		Return(errors.New("an error"))
	err = rt.RunContainer(&Network{}, con)
	if !errors.Is(err, FailedToStartContainer) {
		t.Errorf("Expected FailedToStartContainer but found %v", err)
	}
	dc.
		EXPECT().
		ContainerStop(any, any, any).
		Return(errors.New("an error"))
	err = rt.StopContainer(con)
	if !errors.Is(err, FailedToStopContainer) {
		t.Errorf("Expected FailedToStopContainer but found %v", err)
	}
}

func TestContainerRemove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)
	con := &Container{
		ID:     "testContainerID",
		Name:   "test",
		Image:  Image{Name: "fold/test"},
		Mounts: []Mount{{"foo", "bar"}},
	}

	dc.
		EXPECT().
		ContainerRemove(any, "testContainerID", any)
	err := rt.RemoveContainer(con)
	if err != nil {
		t.Errorf("Expected error to be nil but found %v", err)
	}
}

func TestContainerRemoveFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)
	con := &Container{
		ID:     "testContainerID",
		Name:   "test",
		Image:  Image{Name: "fold/test"},
		Mounts: []Mount{{"foo", "bar"}},
	}

	dc.
		EXPECT().
		ContainerRemove(any, "testContainerID", any).
		Return(errors.New("an error"))
	err := rt.RemoveContainer(con)
	if !errors.Is(err, FailedToRemoveContainer) {
		t.Errorf("Expected FailedToRemoveContainer but found %v", err)
	}
}

func TestContainerJoinAndLeaveNetwork(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)
	net := &Network{Name: "testNet", ID: "testNetID"}
	con := &Container{
		Name:   "testCon",
		ID:     "testConID",
		Image:  Image{Name: "fold/test"},
		Mounts: []Mount{{"foo", "bar"}},
	}

	// Happy
	dc.
		EXPECT().
		NetworkConnect(any, "testNetID", "testConID", any)
	rt.AddToNetwork(net, con)

	dc.
		EXPECT().
		NetworkDisconnect(any, "testNetID", "testConID", false)
	rt.RemoveFromNetwork(net, con)
}

func TestContainerJoinAndLeaveNetworkFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)
	net := &Network{Name: "testNet", ID: "testNetID"}
	con := &Container{
		Name:   "testCon",
		ID:     "testConID",
		Image:  Image{Name: "fold/test"},
		Mounts: []Mount{{"foo", "bar"}},
	}

	dc.
		EXPECT().
		NetworkConnect(any, "testNetID", "testConID", any).
		Return(errors.New("an error"))
	err := rt.AddToNetwork(net, con)
	if !errors.Is(err, FailedToJoinNetwork) {
		t.Errorf("Expected FailedToJoinNetwork but got %v", err)
	}

	dc.
		EXPECT().
		NetworkDisconnect(any, "testNetID", "testConID", false).
		Return(errors.New("an error"))
	err = rt.RemoveFromNetwork(net, con)
	if !errors.Is(err, FailedToLeaveNetwork) {
		t.Errorf("Expected FailedToLeaveNetwork but got %v", err)
	}
}

func TestListAllFoldContainers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)

	dc.
		EXPECT().
		ContainerList(any, any).
		Return([]types.Container{
			containerResponse("a", "/foo", "/bar"),
			containerResponse("b", "/fold.foo", "/bar"),
			containerResponse("c", "/foo", "/fold.bar"),
			containerResponse("d", "/fold.baz"),
			containerResponse("e", "/bar"),
		}, nil)

	cs, err := rt.AllContainers()
	if err != nil {
		t.Errorf("Expected nil but foudn %v", err)
	}
	expectation := []*Container{
		{
			ID:     "b",
			Name:   "fold.foo",
			Image:  Image{Name: "test"},
			Mounts: []Mount{Mount{"src", "dst"}},
		},
		{
			ID:     "c",
			Name:   "fold.bar",
			Image:  Image{Name: "test"},
			Mounts: []Mount{Mount{"src", "dst"}},
		},
		{
			ID:     "d",
			Name:   "fold.baz",
			Image:  Image{Name: "test"},
			Mounts: []Mount{Mount{"src", "dst"}},
		},
	}
	diffContainers(t, expectation, cs)
}

func TestListAllFoldContainersFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)

	dc.
		EXPECT().
		ContainerList(any, any).
		Return([]types.Container{}, errors.New("an error"))

	_, err := rt.AllContainers()
	if !errors.Is(err, DockerEngineError) {
		t.Errorf("Expected DockerEngineError but foudn %v", err)
	}
}

func TestGetContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)

	dc.
		EXPECT().
		ContainerList(any, any).
		Return([]types.Container{
			containerResponse("a", "/foo", "/bar"),
			containerResponse("b", "/fold.foo", "/bar"),
			containerResponse("c", "/foo", "/fold.bar"),
			containerResponse("d", "/fold.baz"),
			containerResponse("e", "/bar"),
		}, nil)

	c, err := rt.GetContainer("foo")
	if err != nil {
		t.Errorf("Expected nil but foudn %v", err)
	}
	expectation := &Container{
		ID: "b", Name: "fold.foo", Image: Image{Name: "test"}, Mounts: []Mount{Mount{"src", "dst"}},
	}
	diffContainers(t, expectation, c)
}

func mockRuntime(dc DockerClient) *ContainerRuntime {
	return &ContainerRuntime{
		cli:    dc,
		ctx:    context.Background(),
		logger: logging.NewTestLogger(),
		out:    os.Stdout,
	}
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
		cmpopts.IgnoreUnexported(Container{}),
	); diff != "" {
		t.Errorf("Expected container lists to be equal(-want +got):\n%s", diff)
	}
}
