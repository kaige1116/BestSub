package task

import "github.com/panjf2000/ants/v2"

var pool *ants.Pool

func Init(maxThread int) {
	pool, _ = ants.NewPool(maxThread)
}

func Submit(fn func()) {
	pool.Submit(fn)
}

func Release() {
	pool.Release()
}
