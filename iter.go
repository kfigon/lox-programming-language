package main


type iter[T any] struct {
	arr []T
	idx int
}

func toIter[T any](arr []T) *iter[T] {
	return &iter[T]{arr, 0}
}

func (i *iter[T]) current() (T, bool) {
	if i.idx >= len(i.arr) {
		var out T
		return out, false
	}
	return i.arr[i.idx], true
}

func (i *iter[T]) peek() (T, bool) {
	if i.idx+1 >= len(i.arr) {
		var out T
		return out, false
	}
	return i.arr[i.idx+1], true
}

func (i *iter[T]) consume() {
	i.idx++
}