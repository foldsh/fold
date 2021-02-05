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
	dcw := mockRuntime(dc)
	net := &Network{Name: "test", rt: dcw}

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
	net.CreateIfNotExists()
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
	net.RemoveIfExists()
}

func TestNetworkCreateFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	dcw := mockRuntime(dc)
	net := &Network{Name: "test", rt: dcw}
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
	err := net.CreateIfNotExists()
	if !errors.Is(err, FailedToCreateNetwork) {
		t.Errorf("Expected FailedToCreateNetwork error but found %v", err)
	}
}

func TestNetworkThatExistsShouldNotBeRecreated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dc := NewMockDockerClient(ctrl)
	dcw := mockRuntime(dc)
	net := &Network{Name: "test", rt: dcw}
	dc.
		EXPECT().
		NetworkList(any, any).
		Return([]types.NetworkResource{
			{Name: "test", ID: "testID"},
			{Name: "foo", ID: "fooID"},
			{Name: "bar", ID: "barID"},
		}, nil)
	err := net.CreateIfNotExists()
	if err != nil {
		t.Errorf("Expected nil but found %v", err)
	}
	if net.ID != "testID" {
		t.Errorf("Expected to find ID of existing network")
	}
}
