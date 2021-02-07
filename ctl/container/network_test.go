package container

import (
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	gomock "github.com/golang/mock/gomock"
)

func TestNetworkCreateAndDestroy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)
	net := &Network{Name: "test"}

	dc.
		EXPECT().
		NetworkList(any, any).
		Return([]types.NetworkResource{
			{Name: "foo", ID: "fooID"},
			{Name: "bar", ID: "barID"},
		}, nil)
	dc.
		EXPECT().
		NetworkCreate(any, "test", any).
		Return(types.NetworkCreateResponse{ID: "testNetID"}, nil)
	rt.NetworkExists(net)
	rt.CreateNetwork(net)
	if net.ID != "testNetID" {
		t.Errorf("After creating expected ID to be 'testNetID' but found %s", net.ID)
	}

	dc.
		EXPECT().
		NetworkList(any, any).
		Return([]types.NetworkResource{
			{Name: "foo", ID: "fooID"},
			{Name: "bar", ID: "barID"},
			{Name: "test", ID: "testNetID"},
		}, nil)
	rt.NetworkExists(net)
	dc.EXPECT().NetworkRemove(any, "testNetID")
	rt.RemoveNetwork(net)
}

func TestNetworkCreateFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)
	net := &Network{Name: "test"}
	dc.
		EXPECT().
		NetworkCreate(any, "test", any).
		Return(types.NetworkCreateResponse{}, errors.New("something went wrong"))
	err := rt.CreateNetwork(net)
	if !errors.Is(err, FailedToCreateNetwork) {
		t.Errorf("Expected FailedToCreateNetwork error but found %v", err)
	}
}

func TestNetworkThatExistsShouldNotBeRecreated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	rt := mockRuntime(dc)
	net := &Network{Name: "test"}
	dc.
		EXPECT().
		NetworkList(any, any).
		Return([]types.NetworkResource{
			{Name: "test", ID: "testID"},
			{Name: "foo", ID: "fooID"},
			{Name: "bar", ID: "barID"},
		}, nil)
	exists, err := rt.NetworkExists(net)
	if err != nil {
		t.Errorf("Expected nil but found %v", err)
	}
	if !exists {
		err := rt.CreateNetwork(net)
		if err != nil {
			t.Errorf("Expected nil but found %v", err)
		}
	}
	if net.ID != "testID" {
		t.Errorf("Expected to find ID of existing network")
	}
}
