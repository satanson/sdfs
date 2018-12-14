package netstg

import (
	"soloos/sdfs/types"
	snettypes "soloos/snet/types"
)

func (p *NetBlockDriver) PRead(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr,
	netBlockIndex int,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int,
	offset, length int) error {
	if uNetBlock.Ptr().DataNodes.Len == 0 {
		return types.ErrBackendListIsEmpty
	}

	var (
		resp snettypes.Response
		err  error
	)

	err = p.dataNodeClient.PRead(uNetBlock.Ptr().DataNodes.Arr[0],
		uNetBlock,
		netBlockIndex,
		uMemBlock,
		memBlockIndex,
		offset, length,
		&resp)
	if err != nil {
		return err
	}

	return nil
}
