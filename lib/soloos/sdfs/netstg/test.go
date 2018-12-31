package netstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/types"
	"soloos/snet"
	snettypes "soloos/snet/types"
	"soloos/util"
	"soloos/util/offheap"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func MakeNetBlockDriversForTest(t *testing.T,
	netBlockDriver *NetBlockDriver,
	offheapDriver *offheap.OffheapDriver,
	snetDriver *snet.NetDriver,
	snetClientDriver *snet.ClientDriver,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
) {
	assert.NoError(t, netBlockDriver.Init(offheapDriver,
		snetDriver, snetClientDriver,
		nameNodeClient, dataNodeClient,
		nil,
	))
}

func MakeDriversForTest(t *testing.T,
	snetDriver *snet.NetDriver,
	snetClientDriver *snet.ClientDriver,
	nameNodeSRPCServerAddr string,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
	netBlockDriver *NetBlockDriver,
) {
	var (
		offheapDriver = &offheap.DefaultOffheapDriver
		nameNodePeer  snettypes.PeerUintptr
	)

	assert.NoError(t, snetDriver.Init(offheapDriver))
	assert.NoError(t, snetClientDriver.Init(offheapDriver))

	nameNodePeer, _ = snetDriver.MustGetPeer(nil, nameNodeSRPCServerAddr, types.DefaultSDFSRPCProtocol)
	assert.NoError(t, nameNodeClient.Init(snetClientDriver, nameNodePeer))
	assert.NoError(t, dataNodeClient.Init(snetClientDriver))
	MakeNetBlockDriversForTest(t, netBlockDriver, offheapDriver,
		snetDriver, snetClientDriver,
		nameNodeClient, dataNodeClient,
	)
}

func MakeMockServerForTest(t *testing.T,
	snetDriver *snet.NetDriver,
	mockServerAddr string, mockServer *MockServer) {
	assert.NoError(t, mockServer.Init(snetDriver, "tcp", mockServerAddr))
	go func() {
		util.AssertErrIsNil(mockServer.Serve())
	}()
	time.Sleep(time.Millisecond * 300)
}

func MakeDriversWithMockServerForTest(t *testing.T,
	snetDriver *snet.NetDriver,
	snetClientDriver *snet.ClientDriver,
	mockServerAddr string,
	mockServer *MockServer,
	nameNodeClient *api.NameNodeClient,
	dataNodeClient *api.DataNodeClient,
	netBlockDriver *NetBlockDriver,
) {
	MakeDriversForTest(t, snetDriver, snetClientDriver, mockServerAddr, nameNodeClient, dataNodeClient, netBlockDriver)
	MakeMockServerForTest(t, snetDriver, mockServerAddr, mockServer)
}
