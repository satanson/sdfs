package types

import (
	"sync"
	"unsafe"
)

const (
	NetINodeIDBytesNum = 64
	NetINodeIDSize     = int(unsafe.Sizeof([NetINodeIDBytesNum]byte{}))
	NetINodeStructSize = unsafe.Sizeof(NetINode{})
)

type NetINodeID = [NetINodeIDBytesNum]byte
type NetINodeUintptr uintptr

func (u NetINodeUintptr) Ptr() *NetINode { return (*NetINode)(unsafe.Pointer(u)) }

type NetINode struct {
	ID               NetINodeID   `db:"netinode_id"`
	Size             int64        `db:"netinode_size"`
	NetBlockCap      int          `db:"netblock_cap"`
	MemBlockCap      int          `db:"memblock_cap"`
	MetaDataMutex    sync.RWMutex `db:"-"`
	IsMetaDataInited bool         `db:"-"`
}

func (p *NetINode) IDStr() string { return string(p.ID[:]) }