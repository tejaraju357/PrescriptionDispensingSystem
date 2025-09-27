package cache

import (
	"context"
	"sync"
	"time"
)

var ctx = context.Background()

var localLocks sync.Map

func getLocalLock(key string) *sync.Mutex {
	val, _ := localLocks.LoadOrStore(key, &sync.Mutex{})
	return val.(*sync.Mutex)
}

func AcquireLock(key string, ttl time.Duration) (bool, error) {
	getLocalLock(key).Lock()

	success, err := Rdb.SetNX(ctx, key, "locked", ttl).Result()
	if err != nil {
		getLocalLock(key).Unlock()
		return false, err
	}

	if !success {
		getLocalLock(key).Unlock()
		return false, nil
	}

	return true, nil
}

func ReleaseLock(key string) error {
	_, err := Rdb.Del(ctx, key).Result()
	getLocalLock(key).Unlock()
	return err
}
