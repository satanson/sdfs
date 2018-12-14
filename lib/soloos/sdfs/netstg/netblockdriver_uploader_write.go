package netstg

import (
	"soloos/sdfs/types"
)

func (p *netBlockDriverUploader) PWrite(uNetINode types.NetINodeUintptr,
	uNetBlock types.NetBlockUintptr,
	uMemBlock types.MemBlockUintptr,
	memBlockIndex int,
	offset, end int) error {

	var (
		isMergeEventHappened    bool
		isMergeWriteMaskSuccess bool = false
		pMemBlock                    = uMemBlock.Ptr()
	)

	for isMergeWriteMaskSuccess == false {
		if pMemBlock.UploadJob.IsUploadPolicyPrepared == false {
			p.prepareUploadMemBlockJob(&pMemBlock.UploadJob,
				uNetINode, uNetBlock, uMemBlock, memBlockIndex)
		}

		pMemBlock.UploadJob.UploadPolicyMutex.Lock()
		isMergeEventHappened, isMergeWriteMaskSuccess =
			pMemBlock.UploadJob.UploadMaskWaiting.Ptr().MergeIncludeNeighbour(offset, end)
		pMemBlock.UploadJob.UploadPolicyMutex.Unlock()

		if isMergeWriteMaskSuccess {
			if isMergeEventHappened == false {
				pMemBlock.UploadJob.UNetINode.Ptr().SyncDataSig.Add(1)
				pMemBlock.UploadJob.SyncDataSig.Add(1)
				p.uploadMemBlockJobChan <- pMemBlock.GetUploadMemBlockJobUintptr()
			}
		}

		if isMergeWriteMaskSuccess == false {
			pMemBlock.UploadJob.SyncDataSig.Wait()
		}
	}

	return nil
}
