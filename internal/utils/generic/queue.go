package generic

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Queue[T Integer] struct {
	Data []T
	Ptr  int
	Full bool
}

func NewQueue[T Integer](capacity int) *Queue[T] {
	if capacity <= 0 {
		panic("queue capacity must be positive")
	}
	return &Queue[T]{
		Data: make([]T, 0, capacity),
		Ptr:  0,
		Full: false,
	}
}

func (q *Queue[T]) Update(value T) {
	if q.Full {
		q.Data[q.Ptr] = value
		q.Ptr = (q.Ptr + 1) % len(q.Data)
	} else {
		q.Data = append(q.Data, value)
		if len(q.Data) == cap(q.Data) {
			q.Full = true
		}
	}
}

func (q *Queue[T]) GetAll() []T {
	if len(q.Data) == 0 {
		return nil
	}

	result := make([]T, len(q.Data))

	if q.Full {
		tailLen := len(q.Data) - q.Ptr
		copy(result, q.Data[q.Ptr:])
		copy(result[tailLen:], q.Data[:q.Ptr])
	} else {
		copy(result, q.Data)
	}

	return result
}

func (q *Queue[T]) Clear() {
	q.Data = q.Data[:0]
	q.Ptr = 0
	q.Full = false
}

func (q *Queue[T]) Average() T {
	if len(q.Data) == 0 {
		return 0
	}

	var sum int64
	for _, value := range q.Data {
		sum += int64(value)
	}

	return T(sum / int64(len(q.Data)))
}
