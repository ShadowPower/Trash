package main

type LinkNode struct {
	prev *LinkNode
	next *LinkNode
	element *GameObject
}

func NewLinkNode(element *GameObject, prev *LinkNode, next *LinkNode) *LinkNode {
	return &LinkNode {
		prev: prev,
		next: next,
		element: element,
	}
}