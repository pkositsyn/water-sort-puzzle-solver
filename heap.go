package watersortpuzzle

import "container/heap"

type distanceHeapElem struct {
	distance     int
	realDistance int
	elem         State
}

type distanceHeap struct {
	heap        []*distanceHeapElem
	elemToIndex map[*distanceHeapElem]int
}

var _ heap.Interface = (*distanceHeap)(nil)

func newDistanceHeap() *distanceHeap {
	return &distanceHeap{heap: make([]*distanceHeapElem, 0), elemToIndex: make(map[*distanceHeapElem]int)}
}

func (d *distanceHeap) Fix(elem *distanceHeapElem) {
	heap.Fix(d, d.elemToIndex[elem])
}

func (d *distanceHeap) Len() int {
	return len(d.heap)
}

func (d *distanceHeap) Less(i, j int) bool {
	return d.heap[i].distance < d.heap[j].distance
}

func (d *distanceHeap) Swap(i, j int) {
	d.elemToIndex[d.heap[i]] = j
	d.elemToIndex[d.heap[j]] = i
	d.heap[i], d.heap[j] = d.heap[j], d.heap[i]
}

func (d *distanceHeap) Push(x interface{}) {
	newElem := x.(*distanceHeapElem)
	d.elemToIndex[newElem] = len(d.heap)
	d.heap = append(d.heap, newElem)
}

func (d *distanceHeap) Pop() interface{} {
	top := d.heap[len(d.heap)-1]
	delete(d.elemToIndex, top)
	d.heap = d.heap[:len(d.heap)-1]
	return top
}
