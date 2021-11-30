package main

type GameObject struct {
	id     int
	x      int
	y      int
	parent *TreeNode
}

func NewGameObject(x int, y int, id int) *GameObject {
	return &GameObject{
		id: id,
		x:  x,
		y:  y,
	}
}
