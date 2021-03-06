package memstg

import (
	"soloos/sdfs/api"
	"soloos/sdfs/netstg"
	"soloos/snet"
	"soloos/util"
	"soloos/util/offheap"
)

func MakeMemBlockDriversForTest(memBlockDriver *MemBlockDriver, offheapDriver *offheap.OffheapDriver,
	blockChunkSize int, blockChunksLimit int32) {
	memBlockDriverOptions := MemBlockDriverOptions{
		[]MemBlockPoolOptions{
			MemBlockPoolOptions{
				blockChunkSize, blockChunksLimit,
			},
		},
	}
	util.AssertErrIsNil(memBlockDriver.Init(offheapDriver, memBlockDriverOptions))
}

func MakeDriversForTest(snetDriver *snet.NetDriver, snetClientDriver *snet.ClientDriver,
	nameNodeSRPCServerAddr string,
	memBlockDriver *MemBlockDriver,
	netBlockDriver *netstg.NetBlockDriver,
	netINodeDriver *NetINodeDriver,
	blockChunkSize int, blockChunksLimit int32) {
	var (
		offheapDriver  = &offheap.DefaultOffheapDriver
		nameNodeClient api.NameNodeClient
		dataNodeClient api.DataNodeClient
	)

	netstg.MakeDriversForTest(snetDriver, snetClientDriver,
		nameNodeSRPCServerAddr,
		&nameNodeClient, &dataNodeClient,
		netBlockDriver,
	)

	MakeMemBlockDriversForTest(memBlockDriver, offheapDriver, blockChunkSize, blockChunksLimit)

	util.AssertErrIsNil(netINodeDriver.Init(offheapDriver, netBlockDriver, memBlockDriver, &nameNodeClient, nil, nil))
}

func MakeDriversWithMockServerForTest(mockServerAddr string,
	mockServer *netstg.MockServer,
	snetDriver *snet.NetDriver,
	netBlockDriver *netstg.NetBlockDriver,
	memBlockDriver *MemBlockDriver,
	netINodeDriver *NetINodeDriver,
	blockChunkSize int, blockChunksLimit int32) {
	var (
		offheapDriver    = &offheap.DefaultOffheapDriver
		snetClientDriver snet.ClientDriver
		nameNodeClient   api.NameNodeClient
		dataNodeClient   api.DataNodeClient
	)

	netstg.MakeDriversWithMockServerForTest(snetDriver, &snetClientDriver,
		mockServerAddr, mockServer,
		&nameNodeClient, &dataNodeClient,
		netBlockDriver,
	)

	MakeMemBlockDriversForTest(memBlockDriver, offheapDriver, blockChunkSize, blockChunksLimit)

	util.AssertErrIsNil(netINodeDriver.Init(offheapDriver, netBlockDriver, memBlockDriver, &nameNodeClient, nil, nil))
}
