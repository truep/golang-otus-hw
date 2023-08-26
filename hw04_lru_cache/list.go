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
	}
	switch l.Front() == nil {
	case true:
		l.base.Prev, l.base.Next = li, li
	default:
		l.Front().Prev = li
		li.Next = l.Front()
		l.base.Prev = li
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
		l.base.Prev, l.base.Next = li, li
	default:
		l.Back().Next = li
		li.Prev = l.Back()
		l.base.Next = li
	}
	l.len++
	return li
}

func (l *list) Remove(i *ListItem) {
	switch {
	case l.Front() == l.Back():
		l.base.Next, l.base.Prev = nil, nil
	case i == l.Front():
		i.Next.Prev = nil
		l.base.Prev = i.Next
	case i == l.Back():
		i.Prev.Next = nil
		l.base.Next = i.Prev
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
		l.base.Next = i.Prev
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	i.Prev = nil
	i.Next = l.Front()
	l.Front().Prev = i
	l.base.Prev = i
}
