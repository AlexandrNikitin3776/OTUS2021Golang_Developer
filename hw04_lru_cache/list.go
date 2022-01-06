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
	firstItem *ListItem
	lastItem  *ListItem
	size      int
}

func NewList() List {
	return new(list)
}

func (l list) Len() int {
	return l.size
}

func (l list) Front() *ListItem {
	return l.firstItem
}

func (l list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.size == 0 {
		l.firstItem = newItem
		l.lastItem = newItem
	} else {
		newItem.Next = l.firstItem
		l.firstItem.Prev = newItem
		l.firstItem = newItem
	}

	l.size++
	return l.firstItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  nil,
	}

	if l.size == 0 {
		l.firstItem = newItem
		l.lastItem = newItem
	} else {
		newItem.Prev = l.lastItem
		l.lastItem.Next = newItem
		l.lastItem = newItem
	}

	l.size++
	return l.lastItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.firstItem = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.lastItem = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.size--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.firstItem == i {
		return
	} else if l.lastItem == i {
		i.Prev.Next = i.Next
		l.lastItem = i.Prev
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	i.Next = l.firstItem
	i.Prev = nil
	l.firstItem.Prev = i
	l.firstItem = i
}
