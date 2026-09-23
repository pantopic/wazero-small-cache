package wazero_small_cache

import (
	"context"
	"encoding/hex"
	"log"
	"sync"
	// "sync/atomic"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tidwall/btree"
)

// Name is the name of this host module.
const Name = "pantopic/wazero-small-cache"

var (
	ctxKeyMeta  = Name + `/meta`
	ctxKeyLocal = Name + `/local`
	ctxKeyMutex = Name + `/mutex`
)

type meta struct {
	ptrGlobal uint32
	ptrID     uint32
	ptrKeyCap uint32
	ptrKeyLen uint32
	ptrKey    uint32
	ptrValCap uint32
	ptrValLen uint32
	ptrVal    uint32
}

type hostModule struct {
	sync.RWMutex

	module api.Module
}

type Option func(*hostModule)

func New(opts ...Option) (h *hostModule) {
	h = &hostModule{}
	for _, opt := range opts {
		opt(h)
	}
	return
}

func (h *hostModule) Name() string {
	return Name
}

func (h *hostModule) Stop() {}

// Register instantiates the host module, making it available to all module instances in this runtime
func (h *hostModule) Register(ctx context.Context, r wazero.Runtime) (err error) {
	builder := r.NewHostModuleBuilder(Name)
	register := func(name string, fn func(ctx context.Context, m api.Module, stack []uint64)) {
		builder = builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(fn), nil, nil).Export(name)
	}
	for name, fn := range map[string]any{
		"__small_cache_put": func(mu *sync.RWMutex, m *btree.Map[string, []byte], k string, v []byte) {
			mu.Lock()
			defer mu.Unlock()
			m.Set(k, v)
		},
		"__small_cache_get": func(mu *sync.RWMutex, m *btree.Map[string, []byte], k string) (v []byte) {
			mu.RLock()
			defer mu.RUnlock()
			v, _ = m.Get(k)
			return
		},
		"__small_cache_del": func(mu *sync.RWMutex, m *btree.Map[string, []byte], k string) {
			mu.Lock()
			defer mu.Unlock()
			m.Delete(k)
		},
		"__small_cache_min": func(mu *sync.RWMutex, m *btree.Map[string, []byte]) (k string) {
			mu.RLock()
			defer mu.RUnlock()
			k, _, _ = m.Min()
			return
		},
	} {
		switch fn := fn.(type) {
		case func(mu *sync.RWMutex, m *btree.Map[string, []byte], k string, v []byte):
			register(name, func(ctx context.Context, mod api.Module, stack []uint64) {
				meta := get[*meta](ctx, ctxKeyMeta)
				fn(get[*sync.RWMutex](ctx, ctxKeyMutex), h.getMap(ctx, mod, meta), getKey(mod, meta), getVal(mod, meta))
			})
		case func(mu *sync.RWMutex, m *btree.Map[string, []byte], k string) (v []byte):
			register(name, func(ctx context.Context, mod api.Module, stack []uint64) {
				meta := get[*meta](ctx, ctxKeyMeta)
				b := fn(get[*sync.RWMutex](ctx, ctxKeyMutex), h.getMap(ctx, mod, meta), getKey(mod, meta))
				copy(valBuf(mod, meta)[:len(b)], b)
				writeUint32(mod, meta.ptrValLen, uint32(len(b)))
			})
		case func(mu *sync.RWMutex, m *btree.Map[string, []byte], k string):
			register(name, func(ctx context.Context, mod api.Module, stack []uint64) {
				meta := get[*meta](ctx, ctxKeyMeta)
				fn(get[*sync.RWMutex](ctx, ctxKeyMutex), h.getMap(ctx, mod, meta), getKey(mod, meta))
			})
		case func(mu *sync.RWMutex, m *btree.Map[string, []byte]) string:
			register(name, func(ctx context.Context, mod api.Module, stack []uint64) {
				meta := get[*meta](ctx, ctxKeyMeta)
				k := fn(get[*sync.RWMutex](ctx, ctxKeyMutex), h.getMap(ctx, mod, meta))
				if len(k) == 0 {
					writeUint32(mod, meta.ptrKeyLen, 0)
					return
				}
				b, err := hex.DecodeString(k)
				if err != nil {
					panic("Unable to decode key in min: " + k)
				}
				copy(keyBuf(mod, meta)[:len(b)], b)
				writeUint32(mod, meta.ptrKeyLen, uint32(len(b)))
			})
		default:
			log.Panicf("Method signature implementation missing: %#v", fn)
		}
	}
	h.module, err = builder.Instantiate(ctx)
	return
}

// InitContext retrieves the meta page from the wasm module
func (h *hostModule) InitContext(ctx context.Context, m api.Module) (context.Context, error) {
	stack, err := m.ExportedFunction(`__small_cache`).Call(ctx)
	if err != nil {
		return ctx, err
	}
	meta := &meta{}
	ptr := uint32(stack[0])
	for i, v := range []*uint32{
		&meta.ptrID,
		&meta.ptrKeyCap,
		&meta.ptrKeyLen,
		&meta.ptrKey,
		&meta.ptrValCap,
		&meta.ptrValLen,
		&meta.ptrVal,
	} {
		*v = readUint32(m, ptr+uint32(4*i))
	}
	return context.WithValue(ctx, ctxKeyMeta, meta), nil
}

func (h *hostModule) ContextCopy(dst, src context.Context) context.Context {
	if v := src.Value(ctxKeyMeta); v != nil {
		dst = context.WithValue(dst, ctxKeyMeta, v.(*meta))
		if v := src.Value(ctxKeyMutex); v != nil {
			dst = context.WithValue(dst, ctxKeyMutex, v.(*sync.RWMutex))
			if v := src.Value(ctxKeyLocal); v != nil {
				dst = context.WithValue(dst, ctxKeyLocal, v.(map[uint64]*btree.Map[string, []byte]))
			} else {
				dst = context.WithValue(dst, ctxKeyLocal, make(map[uint64]*btree.Map[string, []byte]))
			}
		} else {
			dst = context.WithValue(dst, ctxKeyMutex, &sync.RWMutex{})
		}
	}
	return dst
}

func (h *hostModule) getMap(ctx context.Context, mod api.Module, meta *meta) *btree.Map[string, []byte] {
	id := readUint64(mod, meta.ptrID)
	m := get[map[uint64]*btree.Map[string, []byte]](ctx, ctxKeyLocal)
	h.RLock()
	_, ok := m[id]
	h.RUnlock()
	if !ok {
		h.Lock()
		if _, ok := m[id]; !ok {
			m[id] = &btree.Map[string, []byte]{}
		}
		h.Unlock()
	}
	return m[id]
}

func getKey(mod api.Module, meta *meta) string {
	b := read(mod, meta.ptrKey, meta.ptrKeyLen, meta.ptrKeyCap)
	return hex.EncodeToString(b)
}

func getVal(mod api.Module, meta *meta) []byte {
	b := read(mod, meta.ptrVal, meta.ptrValLen, meta.ptrValCap)
	return append([]byte(nil), b...)
}

func valBuf(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrVal, 0, meta.ptrValCap)
}

func keyBuf(m api.Module, meta *meta) []byte {
	return read(m, meta.ptrKey, 0, meta.ptrKeyCap)
}

func get[T any](ctx context.Context, key string) T {
	v := ctx.Value(key)
	if v == nil {
		log.Panicf("Context item missing %s", key)
	}
	return v.(T)
}

func id(m api.Module, meta *meta) uint32 {
	return readUint32(m, meta.ptrID)
}

func readUint32(m api.Module, ptr uint32) (val uint32) {
	val, ok := m.Memory().ReadUint32Le(ptr)
	if !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
	return
}

func read(m api.Module, ptrData, ptrLen, ptrMax uint32) (buf []byte) {
	buf, ok := m.Memory().Read(ptrData, readUint32(m, ptrMax))
	if !ok {
		log.Panicf("Memory.Read(%d, %d) out of range", ptrData, ptrLen)
	}
	return buf[:readUint32(m, ptrLen)]
}

func readUint64(m api.Module, ptr uint32) (val uint64) {
	val, ok := m.Memory().ReadUint64Le(ptr)
	if !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
	return
}

func writeUint64(m api.Module, ptr uint32, val uint64) {
	if ok := m.Memory().WriteUint64Le(ptr, val); !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
}

func writeUint32(m api.Module, ptr uint32, val uint32) {
	if ok := m.Memory().WriteUint32Le(ptr, val); !ok {
		log.Panicf("Memory.Read(%d) out of range", ptr)
	}
}
