package generic

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Queue[T Integer] struct {
	data []T
	ptr  int
	full bool
}

func NewQueue[T Integer](capacity int) *Queue[T] {
	if capacity <= 0 {
		panic("queue capacity must be positive")
	}
	return &Queue[T]{
		data: make([]T, 0, capacity),
		ptr:  0,
		full: false,
	}
}

func (q *Queue[T]) Update(value T) {
	if q.full {
		q.data[q.ptr] = value
		q.ptr = (q.ptr + 1) % len(q.data)
	} else {
		q.data = append(q.data, value)
		if len(q.data) == cap(q.data) {
			q.full = true
		}
	}
}

func (q *Queue[T]) GetAll() []T {
	if len(q.data) == 0 {
		return nil
	}

	result := make([]T, len(q.data))

	if q.full {
		tailLen := len(q.data) - q.ptr
		copy(result, q.data[q.ptr:])
		copy(result[tailLen:], q.data[:q.ptr])
	} else {
		copy(result, q.data)
	}

	return result
}

func (q *Queue[T]) Clear() {
	q.data = q.data[:0]
	q.ptr = 0
	q.full = false
}

func (q *Queue[T]) Average() T {
	if len(q.data) == 0 {
		return 0
	}

	var sum int64
	for _, value := range q.data {
		sum += int64(value)
	}

	return T(sum / int64(len(q.data)))
}
