package main

import (
	"encoding/binary"

	"github.com/pantopic/wazero-small-cache/sdk-go"
)

const (
	SMALL_CACHE_ID_TEST_1 = iota
	SMALL_CACHE_ID_TEST_2
)

var (
	testLocalCache1 *small_cache.Local
	testLocalCache2 *small_cache.Local
)

func main() {
	testLocalCache1 = small_cache.NewLocal(SMALL_CACHE_ID_TEST_1)
	testLocalCache2 = small_cache.NewLocal(SMALL_CACHE_ID_TEST_2)
}

//export testLocalPut
func testLocalPut(k, v uint64) {
	testLocalCache1.Put(
		binary.LittleEndian.AppendUint64([]byte{}, k),
		binary.LittleEndian.AppendUint64([]byte{}, v),
	)
}

//export testLocalGet
func testLocalGet(k uint64) uint64 {
	b := testLocalCache1.Get(binary.LittleEndian.AppendUint64([]byte{}, k))
	if len(b) != 8 {
		return 0
	}
	return binary.LittleEndian.Uint64(b)
}

//export testLocalDel
func testLocalDel(k uint64) {
	testLocalCache1.Del(binary.LittleEndian.AppendUint64([]byte{}, k))
}

//export testLocalMin
func testLocalMin() uint64 {
	b := testLocalCache1.Min()
	return binary.LittleEndian.Uint64(b)
}

//export testLocalPut2
func testLocalPut2(k, v uint64) {
	testLocalCache2.Put(
		binary.LittleEndian.AppendUint64([]byte{}, k),
		binary.LittleEndian.AppendUint64([]byte{}, v),
	)
}

//export testLocalGet2
func testLocalGet2(k uint64) uint64 {
	b := testLocalCache2.Get(binary.LittleEndian.AppendUint64([]byte{}, k))
	if len(b) != 8 {
		return 0
	}
	return binary.LittleEndian.Uint64(b)
}

//export testLocalDel2
func testLocalDel2(k uint64) {
	testLocalCache2.Del(binary.LittleEndian.AppendUint64([]byte{}, k))
}

//export testLocalMin2
func testLocalMin2() uint64 {
	b := testLocalCache2.Min()
	return binary.LittleEndian.Uint64(b)
}
