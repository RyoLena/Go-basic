package pkg

import (
	"reflect"
	"sync"
)

func SliceAdd[T any](slice []T, item T) []T {
	return append(slice, item)
}

func SliceRemove[T comparable](slice []T, item T) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if v != item {
			result = append(result, v)
		}
	}
	return result
}

func SliceContains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func SliceUnion[T comparable](slice1 []T, slice2 []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(slice1)+len(slice2))

	for _, v := range slice1 {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	for _, v := range slice2 {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result

}

func SliceIntersection[T comparable](slice1 []T, slice2 []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0)
	for _, v := range slice1 {
		seen[v] = true
	}

	for _, v := range slice2 {
		if seen[v] {
			result = append(result, v)
		}
	}
	return result
}

// SliceMap Map function for slices
func SliceMap[T any, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

// SliceReduce reduce function for slice
func SliceReduce[T any, U any](slice []T, init U, f func(U, T) U) U {
	result := init
	for _, v := range slice {
		result = f(result, v)
	}
	return result
}

// ——————————————————————————map————————————————————————————————————

func MapGet[K comparable, V any](m map[K]V, key K) (V, bool) {
	val, ok := m[key]
	return val, ok
}

// HashMap implementation
type HashMap[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		data: make(map[K]V),
	}
}

func (hm *HashMap[K, V]) Put(key K, value V) {
	hm.mu.Lock()
	defer hm.mu.Unlock()
	hm.data[key] = value
}

func (hm *HashMap[K, V]) Get(key K) (V, bool) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()
	value, ok := hm.data[key]
	return value, ok
}

// ——————————————————————————list————————————————————————————————————

// ArrayList implementation
type ArrayList[T any] struct {
	data []T
	mu   sync.RWMutex
}

func NewArrayList[T any]() *ArrayList[T] {
	return &ArrayList[T]{
		data: make([]T, 0),
	}
}
func (al *ArrayList[T]) Add(item T) {
	al.mu.Lock()
	defer al.mu.Unlock()
	al.data = append(al.data, item)
}
func (al *ArrayList[T]) Get(index int) (T, bool) {
	al.mu.RLock()
	defer al.mu.RUnlock()
	if index < 0 || index >= len(al.data) {
		var zero T
		return zero, false
	}
	return al.data[index], true
}

// ——————————————————————————set————————————————————————————————————

// HashSet implementation
type HashSet[T comparable] struct {
	data map[T]bool
	mu   sync.RWMutex
}

func NewHashSet[T comparable]() *HashSet[T] {
	return &HashSet[T]{
		data: make(map[T]bool),
	}
}
func (hs *HashSet[T]) Add(item T) {
	hs.mu.Lock()
	defer hs.mu.Unlock()
	hs.data[item] = true
}

func (hs *HashSet[T]) Contains(item T) bool {
	hs.mu.RLock()
	defer hs.mu.RUnlock()
	_, ok := hs.data[item]
	return ok
}

// ——————————————————————————queue————————————————————————————————————

// Queue implementation
type Queue[T any] struct {
	data []T
	mu   sync.Mutex
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		data: make([]T, 0),
	}
}

func (q *Queue[T]) Enqueue(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data = append(q.data, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}
	item := q.data[0]
	q.data = q.data[1:]
	return item, true
}

// ——————————————————————————简易拷贝————————————————————————————————————
// BeanCopy BeanCopier basic implementation, no error handling
func BeanCopy(dest, src interface{}) {
	destVal := reflect.ValueOf(dest).Elem()
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	for i := 0; i < destVal.NumField(); i++ {
		destField := destVal.Type().Field(i)
		srcField := srcVal.FieldByName(destField.Name)
		if srcField.IsValid() && srcField.Type() == destField.Type {
			destVal.Field(i).Set(srcField)
		}
	}
}

// ConcurrentQueue 并发队列 (示例：使用 channel)
type ConcurrentQueue[T any] struct {
	data chan T
}

// NewConcurrentQueue 创建一个带有缓冲大小的 concurrentQueue
func NewConcurrentQueue[T any](cap int) *ConcurrentQueue[T] {
	return &ConcurrentQueue[T]{
		data: make(chan T, cap),
	}
}
func (cq *ConcurrentQueue[T]) Enqueue(item T) {
	cq.data <- item
}
func (cq *ConcurrentQueue[T]) Dequeue() (T, bool) {
	item, ok := <-cq.data
	return item, ok
}

// WorkerPool 简单的协程池
type WorkerPool struct {
	jobQueue    chan func()
	wg          sync.WaitGroup
	workerCount int
}

func NewWorkerPool(workerCount int) *WorkerPool {
	return &WorkerPool{
		jobQueue:    make(chan func()),
		workerCount: workerCount,
	}
}
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerCount; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for job := range wp.jobQueue {
				job()
			}
		}()
	}
}
func (wp *WorkerPool) Submit(job func()) {
	wp.jobQueue <- job
}
func (wp *WorkerPool) Stop() {
	close(wp.jobQueue)
	wp.wg.Wait()
}
