package small_cache

import (
	"unsafe"
)

var (
	id     uint64
	keyCap uint32 = 256
	keyLen uint32
	key           = make([]byte, keyCap)
	valCap uint32 = 4 << 10 // 4 KiB
	valLen uint32
	val    = make([]byte, valCap)
	meta   = make([]uint32, 8)
)

//export __small_cache
func __small_cache() (res uint32) {
	for i, p := range []unsafe.Pointer{
		unsafe.Pointer(&id),
		unsafe.Pointer(&keyCap),
		unsafe.Pointer(&keyLen),
		unsafe.Pointer(&key[0]),
		unsafe.Pointer(&valCap),
		unsafe.Pointer(&valLen),
		unsafe.Pointer(&val[0]),
	} {
		meta[i] = uint32(uintptr(p))
	}
	return uint32(uintptr(unsafe.Pointer(&meta[0])))
}

//go:wasm-module pantopic/wazero-atomic
//export __small_cache_put
func put()

//go:wasm-module pantopic/wazero-atomic
//export __small_cache_get
func get()

//go:wasm-module pantopic/wazero-atomic
//export __small_cache_del
func del()

// Fix for lint rule `unusedfunc`
var _ = __small_cache
