package task

import "github.com/panjf2000/ants/v2"

var pool *ants.Pool
var thread int

func Init(maxThread int) {
	pool, _ = ants.NewPool(maxThread)
	thread = maxThread
}

func Submit(fn func()) {
	pool.Submit(fn)
}

func Release() {
	pool.Release()
}

func MaxThread() int {
	return thread
}
