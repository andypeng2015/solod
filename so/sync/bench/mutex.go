package main

import (
	"solod.dev/so/conc"
	"solod.dev/so/mem"
	"solod.dev/so/sync"
	"solod.dev/so/testing"
)

// numWorkers is the number of threads contending for the mutex
// in the contended benchmark.
const numWorkers = 8

// numLoops is the number of Lock/Unlock rounds each worker performs per
// benchmark iteration. It is large enough to amortize the pool submission
// and thread-wakeup overhead so the measurement reflects lock contention.
const numLoops = 1000

func BenchmarkMutexUncontended_So(b *testing.B) {
	// Measures Lock/Unlock on a mutex that is never contended,
	// i.e. the fast path of the primitive.
	var mu sync.Mutex
	mu.Init()
	defer mu.Free()
	for b.Loop() {
		mu.Lock()
		mu.Unlock()
	}
}

func BenchmarkMutexTryLock_So(b *testing.B) {
	// Measures TryLock/Unlock on an uncontended mutex.
	// TryLock always succeeds here since nothing else holds the lock.
	var mu sync.Mutex
	mu.Init()
	defer mu.Free()
	for b.Loop() {
		if mu.TryLock() {
			mu.Unlock()
		}
	}
}

func BenchmarkMutexContended_So(b *testing.B) {
	// Measures Lock/Unlock under contention:
	// numWorkers threads each hammer the same mutex.
	var mu sync.Mutex
	mu.Init()
	defer mu.Free()

	opts := conc.PoolOpts{NumThreads: numWorkers}
	p := conc.NewPool(mem.System, opts)
	defer p.Free()

	for b.Loop() {
		for range numWorkers {
			p.Go(hammerMutex, &mu)
		}
		p.Wait()
	}
}

// hammerMutex locks and unlocks the shared mutex numLoops times.
func hammerMutex(arg any) {
	mu := arg.(*sync.Mutex)
	for range numLoops {
		mu.Lock()
		mu.Unlock()
	}
}
