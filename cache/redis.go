package cache

import (
	"log"
	"time"

	"github.com/go-redis/redis/v7"
)

type redisClient struct {
	client *redis.Client
}

// Nil defines redis returned nil value error.
const Nil = redis.Nil

var singleton redisClient

// Initialize init package.
func Initialize(addr, password string) {
	singleton.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	err := Redis().Ping().Err()
	if err != nil {
		log.Panicf("connect to redis(%v) failed: %v", addr, err)
	}
}

// InitializeSentinel init package.
func InitializeSentinel(addr, password, masterName, sentinelPassword string) {
	singleton.client = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       masterName,
		SentinelAddrs:    []string{addr},
		SentinelPassword: sentinelPassword,
		Password:         password,
	})
	err := Redis().Ping().Err()
	if err != nil {
		log.Panicf("connect to redis(%v) failed: %v", addr, err)
	}
}

// Redis returns redis client. It's safe to concurrent use.
func Redis() *redis.Client {
	if singleton.client == nil {
		panic("redis client is not created")
	}
	return singleton.client
}

// Lock implements lock using redis.
type Lock struct {
	client *redis.Client
	key    string
}

// NewLock creates lock.
func NewLock(c *redis.Client, key string) *Lock {
	return &Lock{client: c, key: key + ":lock"}
}

// Lock blocks and wait until locks.
func (l *Lock) Lock() {
	for {
		if l.TryLock() {
			return
		}

		time.Sleep(10 * time.Millisecond)
	}
}

// TryLock try to lock and returns result.
func (l *Lock) TryLock() bool {
	return l.client.SetNX(l.key, "1", 30*time.Second).Val()
}

// Unlock unlocks.
func (l *Lock) Unlock() {
	l.client.Del(l.key)
}

// Locked returns whether locked.
func (l *Lock) Locked() bool {
	return l.client.Get(l.key).Val() == "1"
}
