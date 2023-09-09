package structure

type Node[T any] struct {
	Value    T
	Children []Node[T]
}
