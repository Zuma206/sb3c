package utils

type List[T any] struct {
	head *ListNode[T]
	tail *ListNode[T]
}

func NewList[T any]() *List[T] {
	return &List[T]{}
}

type ListNode[T any] struct {
	next  *ListNode[T]
	value T
}

func (list *List[T]) PushBack(value T) {
	node := &ListNode[T]{value: value}
	if list.head == nil {
		list.head = node
	}
	if list.tail != nil {
		list.tail.next = node
	}
	list.tail = node
}
