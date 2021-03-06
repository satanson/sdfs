package main

import "C"
import (
	"io"
	"reflect"
	"soloos/log"
	"soloos/sdfs/libsdfs"
	"soloos/sdfs/types"
	"unsafe"
)

//export GoSdfsPappend
func GoSdfsPappend(fdID uint64, buffer unsafe.Pointer, bufferLen, offset int32) (int32, C.int) {
	var (
		fsINode types.FsINode
		fd      = env.Client.FileTableGet(fdID)
		err     error
	)

	fsINode, err = env.Client.MetaStg.DirTreeDriver.GetFsINodeByID(fd.FsINodeID)
	if err != nil {
		return 0, libsdfs.CODE_ERR
	}

	var data = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufferLen),
		Cap:  int(bufferLen),
	}))

	err = env.Client.MemStg.NetINodeDriver.PWriteWithMem(fsINode.UNetINode, data, int64(offset))
	if err != nil {
		return 0, libsdfs.CODE_ERR
	}

	return bufferLen, 0
}

//export GoSdfsAppend
func GoSdfsAppend(fdID uint64, buffer unsafe.Pointer, bufferLen int32) (int32, C.int) {
	var (
		fsINode types.FsINode
		fd      = env.Client.FileTableGet(fdID)
		err     error
	)

	fsINode, err = env.Client.MetaStg.DirTreeDriver.GetFsINodeByID(fd.FsINodeID)
	if err != nil {
		return 0, libsdfs.CODE_ERR
	}

	var data = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufferLen),
		Cap:  int(bufferLen),
	}))
	err = env.Client.MemStg.NetINodeDriver.PWriteWithMem(fsINode.UNetINode, data, fd.AppendPosition)
	if err != nil {
		log.Warn(err)
		return 0, libsdfs.CODE_ERR
	}

	env.Client.FileTableAddAppendPosition(fdID, int64(bufferLen))

	return bufferLen, libsdfs.CODE_OK
}

//export GoSdfsRead
func GoSdfsRead(fdID uint64, buffer unsafe.Pointer, bufferLen int32) (int32, C.int) {
	var (
		fsINode        types.FsINode
		fd             = env.Client.FileTableGet(fdID)
		readDataLength int
		err            error
	)

	fsINode, err = env.Client.MetaStg.DirTreeDriver.GetFsINodeByID(fd.FsINodeID)
	if err != nil {
		return 0, libsdfs.CODE_ERR
	}

	var data = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufferLen),
		Cap:  int(bufferLen),
	}))
	readDataLength, err = env.Client.MemStg.NetINodeDriver.PReadWithMem(fsINode.UNetINode, data, fd.ReadPosition)
	if err != nil && err != io.EOF {
		log.Warn(err, readDataLength)
		return int32(readDataLength), libsdfs.CODE_ERR
	}

	env.Client.FileTableAddReadPosition(fdID, int64(bufferLen))

	return int32(readDataLength), libsdfs.CODE_OK
}

//export GoSdfsPread
func GoSdfsPread(fdID uint64, buffer unsafe.Pointer, bufferLen int32, position int64) (int32, C.int) {
	var (
		fsINode        types.FsINode
		fd             = env.Client.FileTableGet(fdID)
		readDataLength int
		err            error
	)

	var data = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(buffer),
		Len:  int(bufferLen),
		Cap:  int(bufferLen),
	}))
	fsINode, err = env.Client.MetaStg.DirTreeDriver.GetFsINodeByID(fd.FsINodeID)
	if err != nil {
		log.Warn(err)
		return 0, libsdfs.CODE_ERR
	}

	readDataLength, err = env.Client.MemStg.NetINodeDriver.PReadWithMem(fsINode.UNetINode, data, position)
	if err != nil {
		return int32(readDataLength), libsdfs.CODE_ERR
	}

	return int32(readDataLength), libsdfs.CODE_OK
}

//export GoSdfsCloseFile
func GoSdfsCloseFile(fdID uint64) C.int {
	return doFlushINode(fdID)
}

//export GoSdfsFlushFile
func GoSdfsFlushFile(fdID uint64) C.int {
	return doFlushINode(fdID)
}

//export GoSdfsHFlushINode
func GoSdfsHFlushINode(fdID uint64) C.int {
	return doFlushINode(fdID)
}

//export GoSdfsHSyncINode
func GoSdfsHSyncINode(fdID uint64) C.int {
	return doFlushINode(fdID)
}

func doFlushINode(fdID uint64) C.int {
	var (
		fsINode types.FsINode
		fd      = env.Client.FileTableGet(fdID)
		err     error
	)

	fsINode, err = env.Client.MetaStg.DirTreeDriver.GetFsINodeByID(fd.FsINodeID)
	if err != nil {
		return libsdfs.CODE_ERR
	}

	if fsINode.UNetINode != 0 {
		env.Client.MemStg.NetINodeDriver.Flush(fsINode.UNetINode)
	}

	return libsdfs.CODE_OK
}
