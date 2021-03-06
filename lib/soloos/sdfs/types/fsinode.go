package types

import (
	"unsafe"
)

type FsINodeID = int64

type FsINodeUintptr uintptr

func (u FsINodeUintptr) Ptr() *FsINode { return (*FsINode)(unsafe.Pointer(u)) }

type FsINode struct {
	ID         FsINodeID
	ParentID   FsINodeID
	Name       string
	Flag       int
	Permission int
	NetINodeID NetINodeID
	Type       int
	CTime      int64
	MTime      int64
	UNetINode  NetINodeUintptr
}

type FsINodeFileHandler struct {
	FsINodeID      FsINodeID
	AppendPosition int64
	ReadPosition   int64
}

func (p *FsINodeFileHandler) Reset() {
	p.AppendPosition = 0
	p.ReadPosition = 0
}
