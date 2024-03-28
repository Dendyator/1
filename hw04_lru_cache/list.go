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
	next, prev *ListItem
	list       *list
	Value      interface{}
}

type list struct {
	root ListItem
	len  int
}

func NewList() List { return new(list).init() }

func (e *ListItem) Next() *ListItem {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

func (e *ListItem) Prev() *ListItem {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

func (l *list) init() *list {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

func (l *list) Len() int { return l.len }

func (l *list) Front() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

func (l *list) Back() *ListItem {
	if l.len == 0 {
		return nil
	}
	return l.root.prev
}

func (l *list) lazyInit() {
	if l.root.next == nil {
		l.init()
	}
}

func (l *list) insert(e, at *ListItem) *ListItem {
	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
	e.list = l
	l.len++
	return e
}

func (l *list) insertValue(v any, at *ListItem) *ListItem {
	return l.insert(&ListItem{Value: v}, at)
}

func (l *list) remove(e *ListItem) {
	e.prev.next = e.next
	e.next.prev = e.prev
	e.next = nil
	e.prev = nil
	e.list = nil
	l.len--
}

func (l *list) move(e, at *ListItem) {
	if e == at {
		return
	}
	e.prev.next = e.next
	e.next.prev = e.prev

	e.prev = at
	e.next = at.next
	e.prev.next = e
	e.next.prev = e
}

func (l *list) Remove(i *ListItem) {
	if i.list == l {
		l.remove(i)
	}
}

func (l *list) PushFront(v any) *ListItem {
	l.lazyInit()
	return l.insertValue(v, &l.root)
}

func (l *list) PushBack(v any) *ListItem {
	l.lazyInit()
	return l.insertValue(v, l.root.prev)
}

func (l *list) MoveToFront(i *ListItem) {
	if i.list != l || l.root.next == i {
		return
	}
	l.move(i, &l.root)
}
