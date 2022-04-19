package Hytrix

import (
	"sync"
	"time"
)



type numberBucket struct {
	Value float64
}
type Number struct {
	Buckets map[int64]*numberBucket
	Mutex   *sync.RWMutex
}

func (r Number) getCurrentBucket() *numberBucket {
	now := time.Now().Unix()

	var bucket *numberBucket
	var ok bool

	if bucket, ok = r.Buckets[now]; !ok {
		bucket = &numberBucket{}
		r.Buckets[now] = bucket
	}

	return bucket
}

func (r Number) removeOldBuckets() {
	now := time.Now().Unix() - 10
	for key := range r.Buckets {
		if key <= now {
			delete(r.Buckets, key)
		}
	}
}

func (r Number) Increment(i float64) {
	if i == 0 {
		return
	}

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	bucket := r.getCurrentBucket()
	bucket.Value += i

	r.removeOldBuckets()
}

