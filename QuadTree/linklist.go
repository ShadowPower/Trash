package main

type LinkList struct {
	size int
	head *LinkNode
	tail *LinkNode
}

type ForeachLinkListFunction func(gameObject *GameObject)

func NewLinkList() *LinkList {
	return &LinkList {
		size: 0,
	}
}

func (linkList *LinkList) add(gameObject *GameObject) {
	if gameObject != nil {
		linkList.size += 1
		oldHead := linkList.head
		linkList.head = NewLinkNode(gameObject, nil, oldHead)
		if oldHead != nil {
			oldHead.prev = linkList.head
		} else {
			linkList.tail = linkList.head
		}
	}
}

func (linkList *LinkList) remove(gameObject *GameObject, removeAll bool) {
	cur := linkList.head
	for {
		if cur == nil {
			break
		}
		if gameObject == cur.element {
			if cur.prev != nil {
				cur.prev.next = cur.next
			}
			if cur.next != nil {
				cur.next.prev = cur.prev
			}
			if linkList.tail == cur {
				linkList.tail = cur.prev
			}
			linkList.size -= 1
			if !removeAll {
				return
			}
		}
		cur = cur.next
	}
}

func (linkList *LinkList) clear() {
	linkList.size = 0
	linkList.head = nil
}

func (linkList *LinkList) contains(gameObject *GameObject) bool {
	cur := linkList.head
	for {
		if cur == nil {
			break
		}
		if gameObject == cur.element {
			return true
		}
		cur = cur.next
	}
	return false
}

func (linkList *LinkList) merge(anotherLinkList *LinkList) {
	if anotherLinkList.head != nil {
		if linkList.tail != nil {
			linkList.tail.next = anotherLinkList.head
			anotherLinkList.head.prev = linkList.tail

			linkList.tail = anotherLinkList.tail
			linkList.size = linkList.size + anotherLinkList.size
		} else {
			linkList.size = anotherLinkList.size
			linkList.head = anotherLinkList.head
			linkList.tail = anotherLinkList.tail
		}
	}
	// anotherLinkList 不能用了，否则数据会出错
	anotherLinkList.head = nil
	anotherLinkList.tail = nil
	anotherLinkList.size = 0
}

func (linkList *LinkList) toGameObjectArray() []*GameObject {
	result := make([]*GameObject, linkList.size)
	cur := linkList.head
	i := 0
	for {
		if cur == nil {
			break
		}
		result[i] = cur.element
		cur = cur.next
		i++
	}
	return result
}

func (linkList *LinkList) foreach(foreachFunction ForeachLinkListFunction) {
	cur := linkList.head
	for {
		if cur == nil {
			break
		}
		foreachFunction(cur.element)
		cur = cur.next
	}
}