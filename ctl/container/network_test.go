package container_test

import (
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/foldsh/fold/ctl/container"
)

func TestNetworkCreateAndDestroy(t *testing.T) {
	rt, dc, _ := setup()
	net := &container.Network{Name: "test"}

	dc.
		On("NetworkList", mock.Anything, mock.Anything).
		Return([]types.NetworkResource{
			{Name: "foo", ID: "fooID"},
			{Name: "bar", ID: "barID"},
		}, nil)
	dc.
		On("NetworkCreate", mock.Anything, "test", mock.Anything).
		Return(types.NetworkCreateResponse{ID: "testNetID"}, nil)
	rt.NetworkExists(net)
	rt.CreateNetwork(net)
	assert.Equal(t, "testNetID", net.ID)

	dc.
		On("NetworkList", mock.Anything, mock.Anything).
		Return([]types.NetworkResource{
			{Name: "foo", ID: "fooID"},
			{Name: "bar", ID: "barID"},
			{Name: "test", ID: "testNetID"},
		}, nil)
	rt.NetworkExists(net)
	dc.On("NetworkRemove", mock.Anything, "testNetID").Return(nil)
	rt.RemoveNetwork(net)
	dc.AssertExpectations(t)
}

func TestNetworkCreateFailure(t *testing.T) {
	rt, dc, _ := setup()
	net := &container.Network{Name: "test"}
	dc.
		On("NetworkCreate", mock.Anything, "test", mock.Anything).
		Return(types.NetworkCreateResponse{}, errors.New("something went wrong"))
	err := rt.CreateNetwork(net)
	assert.ErrorIs(t, err, container.FailedToCreateNetwork)
	dc.AssertExpectations(t)
}

func TestNetworkThatExistsShouldNotBeRecreated(t *testing.T) {
	rt, dc, _ := setup()
	net := &container.Network{Name: "test"}
	dc.
		On("NetworkList", mock.Anything, mock.Anything).
		Return([]types.NetworkResource{
			{Name: "test", ID: "testID"},
			{Name: "foo", ID: "fooID"},
			{Name: "bar", ID: "barID"},
		}, nil)
	exists, err := rt.NetworkExists(net)
	assert.Nil(t, err)
	if !exists {
		err := rt.CreateNetwork(net)
		assert.Nil(t, err)
		if err != nil {
			t.Errorf("Expected nil but found %v", err)
		}
	}
	assert.Equal(t, "testID", net.ID)
	dc.AssertExpectations(t)
}
