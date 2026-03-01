package small_cache

type Local struct {
	id uint64
}

func NewLocal(id uint64) *Local {
	return &Local{id}
}

func (c *Local) Get(k []byte) []byte {
	id = c.id
	copy(key[:len(k)], k)
	keyLen = uint32(len(k))
	get()
	return val[:valLen]
}

func (c *Local) Put(k, v []byte) {
	id = c.id
	copy(key[:len(k)], k)
	keyLen = uint32(len(k))
	copy(val[:len(v)], v)
	valLen = uint32(len(v))
	put()
}

func (c *Local) Del(k []byte) {
	id = c.id
	copy(key[:len(k)], k)
	keyLen = uint32(len(k))
	del()
}
