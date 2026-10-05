package imagecache

import (
	"bytes"
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/raceresult/gopdf/types"
)

// ImageCache caches images/mask in mem, so that the same image does not have to be decoded all over again
type ImageCache struct {
	items map[uint64]*Item
	mux   sync.RWMutex

	maxSize  int
	currSize int
}

type Processor func([]byte) (*Item, error)

// New creates a new ImageCache with the given maximum size
func New(maxSize int) *ImageCache {
	return &ImageCache{
		items:   make(map[uint64]*Item),
		maxSize: maxSize,
	}
}

// Process accepts the given image bytes. If already know, it does not need to call the processor function again
func (q *ImageCache) Process(bts []byte, processor Processor) (*Item, error) {
	// already known?
	hash := xxhash.Sum64(bts)
	q.mux.RLock()
	item, ok := q.items[hash]
	q.mux.RUnlock()
	if ok {
		return item, nil
	}

	// process
	item, err := processor(bts)
	if err != nil {
		return nil, err
	}
	item.size = item.Data.Len() + item.Mask.Len()

	// add to cache
	q.mux.Lock()
	defer q.mux.Unlock()
	q.items[hash] = item
	q.currSize += item.size

	// reduce cache (remove random items)
	if q.currSize > q.maxSize {
		for k, v := range q.items {
			q.currSize -= v.size
			delete(q.items, k)
			if q.currSize <= q.maxSize {
				break
			}
		}
	}
	return item, nil
}

type Item struct {
	Data       bytes.Buffer
	Mask       bytes.Buffer
	ColorModel types.ColorSpaceFamily
	size       int
}
