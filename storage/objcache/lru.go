// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package objcache

import (
	"context"
	"time"

	"github.com/corestoreio/pkg/storage/lru"
)

// LRUOptions allows to track the cache items either by object count or total
// size in bytes. If tracking by `TrackBySize` gets enabled then `Capacity` must
// have the value in bytes. Default 64MB.
type LRUOptions struct {
	Capacity           int64 // default 5000 objects
	TrackBySize        bool
	TrackByObjectCount bool // default
	LRUCache           *lru.Cache
}

// lruCache is an LRU cache. It is safe for concurrent access.
type lruCache struct {
	opt LRUOptions
}

// NewLRU creates a new LRU Storage. Expirations are not supported. Argument `o`
// can be nil, if so default values get applied.
func NewLRU(o *LRUOptions) NewStorageFn {
	if o == nil {
		o = &LRUOptions{}
	}
	switch {
	case o.TrackBySize && o.Capacity == 0:
		o.Capacity = 1 << 26 // 64MB
	case o.TrackByObjectCount && o.Capacity == 0:
		o.Capacity = 5000 // objects
	case o.TrackBySize:

	case o.TrackByObjectCount:

	default:
		o.TrackByObjectCount = true
		o.Capacity = 5000
	}
	if o.LRUCache == nil {
		o.LRUCache = lru.New(o.Capacity)
	}
	return func() (Storager, error) {
		return lruCache{
			opt: *o,
		}, nil
	}
}

type itemBySize []byte

func (li itemBySize) Size() int { return len(li) }

type itemByCount []byte

func (li itemByCount) Size() int { return 1 }

func (c lruCache) Set(_ context.Context, keys []string, values [][]byte, _ []time.Duration) (err error) {
	for i, key := range keys {
		var v lru.Value = itemByCount(values[i])
		if c.opt.TrackBySize {
			v = itemBySize(values[i])
		}
		c.opt.LRUCache.Set(key, v)
	}
	return nil
}

// Get looks up a key's value from the cache.
func (c lruCache) Get(_ context.Context, keys []string) (values [][]byte, err error) {
	for _, key := range keys {
		itm, ok := c.opt.LRUCache.Get(key)
		if ok {
			if c.opt.TrackByObjectCount {
				values = append(values, []byte(itm.(itemByCount)))
			} else {
				values = append(values, []byte(itm.(itemBySize)))
			}
		} else {
			values = append(values, nil)
		}
	}
	return
}

func (c lruCache) Truncate(_ context.Context) (err error) {
	c.opt.LRUCache.Clear()
	return nil
}

func (c lruCache) Delete(_ context.Context, keys []string) (err error) {
	for _, key := range keys {
		c.opt.LRUCache.Delete(key)
	}
	return nil
}

func (c lruCache) Close() error {
	c.opt.LRUCache.Clear()
	return nil
}
