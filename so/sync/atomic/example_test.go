package atomic_test

import (
	"solod.dev/so/conc"
	"solod.dev/so/mem"
	"solod.dev/so/sync/atomic"
)

func bump(arg any) {
	cnt := arg.(*atomic.Int64)
	cnt.Add(1)
}

func ExampleInt64() {
	// A shared counter incremented by many threads without a mutex.
	var cnt atomic.Int64

	opts := conc.PoolOpts{NumThreads: 4}
	pool := conc.NewPool(mem.System, opts)
	defer pool.Free()
	for range 100 {
		pool.Go(bump, &cnt)
	}
	pool.Wait()

	println(cnt.Load())
	// 100
}

func ExampleBool() {
	// A stop flag polled by a worker and set from another thread.
	var stop atomic.Bool

	if stop.Load() {
		println("stop requested")
	}
	stop.Store(true)
	if stop.Load() {
		println("stop requested")
	}
	// stop requested
}

type config struct {
	addr string
}

func ExamplePointer() {
	// Publishing a new config that readers pick up atomically.
	var cur atomic.Pointer[config]

	old := config{addr: "localhost:8080"}
	cur.Store(&old)

	next := config{addr: "localhost:9090"}
	if cur.CompareAndSwap(&old, &next) {
		println(cur.Load().addr)
	}
	// localhost:9090
}
