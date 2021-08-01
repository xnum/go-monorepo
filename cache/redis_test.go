package cache_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go-monorepo/cache"
)

func TestRedis(t *testing.T) {
	assert := assert.New(t)

	cache.Initialize("localhost:6379", "")
	redis := cache.Redis()
	redis.FlushAll()
	val, err := redis.Ping().Result()
	assert.NoError(err)
	fmt.Println(val)

	err = redis.Set("ABCD", "XXX", 10*time.Second).Err()
	assert.NoError(err)

	val, err = redis.Get("ABCD").Result()
	assert.Equal("XXX", val)
	assert.NoError(err)

	lock := cache.NewLock(redis, "ABC")
	assert.False(lock.Locked())
	assert.True(lock.TryLock())
	assert.True(lock.Locked())
	lock.Unlock()
	assert.False(lock.Locked())

	ok, err := redis.SetNX("ABCDE", "lu", 3*time.Second).Result()
	assert.True(ok)
	assert.NoError(err)
	ok, err = redis.SetNX("ABCDE", "lu", 3*time.Second).Result()
	assert.False(ok)
	assert.NoError(err)

	s, err := redis.Ping().Result()
	assert.NoError(err)
	assert.Equal(s, "PONG")
}

func BenchmarkLock(b *testing.B) {
	cache.Initialize("localhost:6379", "")
	redis := cache.Redis()
	redis.FlushAll()

	lk := cache.NewLock(redis, "TEST2")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		for j := 0; j < 2; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				if lk.TryLock() {
					defer lk.Unlock()
				} else {
					for {
						if !lk.Locked() {
							break
						}
					}
				}
			}()
		}

		wg.Wait()
	}
}
