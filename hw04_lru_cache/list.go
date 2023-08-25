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
	base ListItem
	len  int
}

func NewList() List {
	return new(list)
}

func (l list) Len() int { return l.len }

func (l list) Front() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.base.Prev
}

func (l list) Back() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.base.Next
}

func (l *list) PushFront(v interface{}) *ListItem {
	li := &ListItem{
		Value: v,
		Next:  l.Front(),
	}

	switch l.len {
	case 0:
		l.base.Prev = li
		l.base.Next = li
	default:
		li.Next = l.Front()
		l.Front().Prev = li
		l.base.Prev = li
	}
	l.len++
	return li
}

func (l *list) PushBack(v interface{}) *ListItem {
	li := &ListItem{
		Value: v,
	}

	switch l.len {
	case 0:
		l.base.Next = li
		l.base.Prev = li
	default:
		li.Prev = l.Back()
		l.Back().Next = li
		l.base.Next = li
	}
	l.len++
	return li
}

func (l *list) Remove(i *ListItem) {
	if i != l.Front() {
		i.Prev.Next = i.Next
	} else {
		l.base.Prev = i.Next
	}
	if i != l.Back() {
		i.Next.Prev = i.Prev
	} else {
		l.base.Next = i.Prev
	}
	i.Next = nil
	i.Prev = nil
	i.Value = nil
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.Front() {
		return
	}

	i.Prev.Next = i.Next
	if i != l.Back() {
		i.Next.Prev = i.Prev
	} else {
		l.base.Next = i.Prev
	}
	i.Prev = nil
	i.Next = l.Front()
	l.Front().Prev = i
	l.base.Prev = i
}
