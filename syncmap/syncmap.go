package syncmap

import "sync"

type SyncByteMap struct {
	sync.RWMutex
	internal map[string][]byte
}

func NewSyncByteMap() *SyncByteMap {
	return &SyncByteMap{
		internal: make(map[string][]byte),
	}
}

func (rm *SyncByteMap) Load(key string) (value []byte, ok bool) {
	rm.RLock()
	result, ok := rm.internal[key]
	rm.RUnlock()
	return result, ok
}

func (rm *SyncByteMap) Delete(key string) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}

func (rm *SyncByteMap) Store(key string, value []byte) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}
