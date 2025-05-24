package _map

// Map default interface
type Map[K comparable, V any] interface {
	Put(key K, value V) (error error)
	Get(key K) (value V, error error)
	Delete(key K) (error error)
	Size() int
}
