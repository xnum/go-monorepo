package pool

import (
	"sync"

	"google.golang.org/grpc"
)

type grpcPoolGetter func() *GRPCPool

var grpcPools sync.Map

// GetGRPCPool provides thread safe singleton.
func GetGRPCPool(key string, opt *Options, dialOptions ...grpc.DialOption) *GRPCPool {
	getter := func() grpcPoolGetter {
		if fn, ok := grpcPools.Load(key); ok {
			return fn.(grpcPoolGetter)
		}

		var pool *GRPCPool
		var once sync.Once
		wrapGetter := grpcPoolGetter(func() *GRPCPool {
			once.Do(func() {
				pool = NewGRPCPool(opt, dialOptions...)
			})

			return pool
		})

		f, loaded := grpcPools.LoadOrStore(key, wrapGetter)
		if loaded {
			return f.(grpcPoolGetter)
		}

		return wrapGetter
	}

	return getter()()
}
