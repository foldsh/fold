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
	rt.CreateNetworkIfNotExists(net)
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
	dc.EXPECT().NetworkRemove(any, "testNetID")
	rt.RemoveNetworkIfExists(net)
}

func TestNetworkCreateFailure(t *testing.T) {
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
		Return(types.NetworkCreateResponse{}, errors.New("something went wrong"))
	err := rt.CreateNetworkIfNotExists(net)
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
	err := rt.CreateNetworkIfNotExists(net)
	if err != nil {
		t.Errorf("Expected nil but found %v", err)
	}
	if net.ID != "testID" {
		t.Errorf("Expected to find ID of existing network")
	}
}
