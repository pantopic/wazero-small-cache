package small_cache

type cache struct {
	id uint64
}

type Local cache

func NewLocal(id uint64) *Local {
	return &Local{id}
}

func (c *cache) Get(k []byte) []byte {
	id = c.id
	copy(key[:len(k)], k)
	keyLen = uint32(len(k))
	get()
	return val[:valLen]
}

func (c *cache) Put(k, v []byte) {
	id = c.id
	copy(key[:len(k)], k)
	keyLen = uint32(len(k))
	copy(val[:len(v)], v)
	valLen = uint32(len(v))
	put()
}

func (c *cache) Del(k []byte) {
	id = c.id
	copy(key[:len(k)], k)
	keyLen = uint32(len(k))
	del()
}
