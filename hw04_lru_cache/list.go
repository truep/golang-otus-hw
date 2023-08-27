package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func NewList() List {
	return new(list)
}

func (l list) Len() int { return l.len }

func (l list) Front() *ListItem {
	return l.head
}

func (l list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := &ListItem{
		Value: v,
	}
	switch l.Front() == nil {
	case true:
		l.head, l.tail = li, li
	default:
		l.Front().Prev = li
		li.Next = l.Front()
		l.head = li
	}
	l.len++
	return li
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := &ListItem{
		Value: v,
	}
	switch l.Back() == nil {
	case true:
		l.head, l.tail = li, li
	default:
		l.Back().Next = li
		li.Prev = l.Back()
		l.tail = li
	}
	l.len++
	return li
}

func (l *list) Remove(i *ListItem) {
	switch {
	case l.Front() == l.Back():
		l.tail, l.head = nil, nil
	case i == l.Front():
		i.Next.Prev = nil
		l.head = i.Next
	case i == l.Back():
		i.Prev.Next = nil
		l.tail = i.Prev
	default:
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	i.Next = nil
	i.Prev = nil
	i.Value = nil
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	switch {
	case i == l.Front():
		return
	case i == l.Back():
		i.Prev.Next = nil
		l.tail = i.Prev
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	i.Prev = nil
	i.Next = l.Front()
	l.Front().Prev = i
	l.head = i
}
