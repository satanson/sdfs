package types

import (
	"reflect"
	"soloos/log"
	snettypes "soloos/snet/types"
	"soloos/util/offheap"
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	MemBlockStructSize = unsafe.Sizeof(MemBlock{})

	MemBlockRefuseReleaseForErr = -1
)

const (
	MemBlockUninited = iota
	MemBlockIniteded
	MemBlockReleasable
	MemBlockRelease
)

type MemBlockUintptr uintptr

func (u MemBlockUintptr) Ptr() *MemBlock {
	return (*MemBlock)(unsafe.Pointer(u))
}

type MemBlock struct {
	ID                  PtrBindIndex
	Status              int64 // equals 0 if could be release
	RebaseNetBlockMutex sync.Mutex
	Chunk               offheap.ChunkUintptr
	Bytes               reflect.SliceHeader
	AvailMask           offheap.ChunkMask
	UploadJob           UploadMemBlockJob
}

func (p *MemBlock) Contains(offset, end int) bool {
	return p.AvailMask.Contains(offset, end)
}

func (p *MemBlock) PWriteWithConn(conn *snettypes.Connection, length int, offset int) (isSuccess bool) {
	_, isSuccess = p.AvailMask.MergeIncludeNeighbour(offset, offset+length)
	if isSuccess {
		var err error
		if offset+length > p.Bytes.Cap {
			length = p.Bytes.Cap - offset
		}
		err = conn.ReadAll((*(*[]byte)(unsafe.Pointer(&p.Bytes)))[offset : offset+length])
		if err != nil {
			log.Warn("PWriteWithConn error", err)
			isSuccess = false
		}
	}
	return
}

func (p *MemBlock) PWriteWithMem(data []byte, offset int) (isSuccess bool) {
	_, isSuccess = p.AvailMask.MergeIncludeNeighbour(offset, offset+len(data))
	if isSuccess {
		copy((*(*[]byte)(unsafe.Pointer(&p.Bytes)))[offset:], data)
	}
	return
}

func (p *MemBlock) PReadWithConn(conn *snettypes.Connection, length int, offset int) error {
	var err error
	err = conn.WriteAll((*(*[]byte)(unsafe.Pointer(&p.Bytes)))[offset : offset+length])
	if err != nil {
		return err
	}
	return nil
}

func (p *MemBlock) PReadWithMem(data []byte, offset int) {
	copy(data, (*(*[]byte)(unsafe.Pointer(&p.Bytes)))[offset:])
}

func (p *MemBlock) GetUploadMemBlockJobUintptr() UploadMemBlockJobUintptr {
	return UploadMemBlockJobUintptr(unsafe.Pointer(p)) + UploadMemBlockJobUintptr(unsafe.Offsetof(p.UploadJob))
}

func (p *MemBlock) BytesSlice() *[]byte {
	return (*[]byte)(unsafe.Pointer(&p.Bytes))
}

func (p *MemBlock) Reset() {
	p.Status = MemBlockUninited
	p.AvailMask.Reset()
	p.UploadJob.Reset()
}

func (p *MemBlock) CompleteInit() {
	p.Status = MemBlockIniteded
}

func (p *MemBlock) IsInited() bool {
	return p.Status > MemBlockUninited
}

func (p *MemBlock) SetReleasable() {
	p.Status = MemBlockReleasable
}

func (p *MemBlock) EnsureRelease() bool {
	return atomic.CompareAndSwapInt64(&p.Status, MemBlockReleasable, MemBlockRelease)
}
