package main

const maxObjects = 4

type TreeNode struct {
	depth       int
	x           int
	y           int
	w           int
	h           int
	objects     *LinkList
	topLeft     *TreeNode
	topRight    *TreeNode
	bottomLeft  *TreeNode
	bottomRight *TreeNode
	parent      *TreeNode
}

func NewTreeNode(depth int, x int, y int, w int, h int, objects *LinkList, parent *TreeNode) *TreeNode {
	return &TreeNode{
		depth:   depth,
		x:       x,
		y:       y,
		w:       w,
		h:       h,
		objects: objects,
		parent:  parent,
	}
}

func (node *TreeNode) addGameObject(gameObject *GameObject) {
	node.objects.add(gameObject)
}

func (node *TreeNode) removeGameObject(gameObject *GameObject) {
	node.objects.remove(gameObject, false)
}

func (node *TreeNode) update() {
	// 可以继续细分的条件
	if node.depth < 20 && node.w > 2 && node.h > 2 {
		if node.objects != nil && node.objects.size > maxObjects {
			var leftW = node.w / 2
			var topH = node.h / 2
			var rightX = node.x + leftW + 1
			var rightW = node.w - leftW - 1
			var bottomY = node.y + topH + 1
			var bottomH = node.h - topH - 1

			node.topLeft = NewTreeNode(node.depth+1, node.x, node.y, leftW, topH, NewLinkList(), node)
			node.topRight = NewTreeNode(node.depth+1, rightX, node.y, rightW, topH, NewLinkList(), node)
			node.bottomLeft = NewTreeNode(node.depth+1, node.x, bottomY, leftW, bottomH, NewLinkList(), node)
			node.bottomRight = NewTreeNode(node.depth+1, rightX, bottomY, rightW, bottomH, NewLinkList(), node)

			cur := node.objects.head
			for {
				if cur == nil {
					break
				}
				gameObject := cur.element
				var targetArea *TreeNode
				switch {
				case gameObject.x < rightX && gameObject.y < bottomY:
					targetArea = node.topLeft
				case gameObject.x >= rightX && gameObject.y < bottomY:
					targetArea = node.topRight
				case gameObject.x < rightX && gameObject.y >= bottomY:
					targetArea = node.bottomLeft
				case gameObject.x >= rightX && gameObject.y >= bottomY:
					targetArea = node.bottomRight
				}
				if targetArea != nil {
					gameObject.parent = targetArea
					targetArea.objects.add(gameObject)
				}
				cur = cur.next
			}

			node.objects = NewLinkList()

			node.topLeft.update()
			node.topRight.update()
			node.bottomLeft.update()
			node.bottomRight.update()

			return
		}
	}
	// 更新不满足条件的 gameObject 的 parent
	node.objects.foreach(func(gameObject *GameObject) {
		gameObject.parent = node
	})
}

func (node *TreeNode) optimize() int {
	if node.topLeft == nil && node.topRight == nil && node.bottomLeft == nil && node.bottomRight == nil {
		// 根节点，返回当前节点对象个数
		return node.objects.size
	} else {
		// 非根节点
		var childObjectNum = node.topLeft.optimize() + node.topRight.optimize() +
			node.bottomLeft.optimize() + node.bottomRight.optimize()
		if childObjectNum <= maxObjects {
			if childObjectNum > 0 {
				// 子节点数据移到父节点
				node.objects.merge(node.topLeft.objects)
				node.objects.merge(node.topRight.objects)
				node.objects.merge(node.bottomLeft.objects)
				node.objects.merge(node.bottomRight.objects)
				node.objects.foreach(func(gameObject *GameObject) {
					gameObject.parent = node
				})
			}

			// 回收内存
			node.topLeft = nil
			node.topRight = nil
			node.bottomLeft = nil
			node.bottomRight = nil
		}
		return childObjectNum
	}
}

func (node *TreeNode) checkOverlap(radius int, xCenter int, yCenter int) bool {
	x1 := node.x
	x2 := node.x + node.w
	y1 := node.y
	y2 := node.y + node.h

	var dx, dy int
	if x1 > xCenter {
		dx = x1 - xCenter
	} else if xCenter > x2 {
		dx = xCenter - x2
	} else {
		dx = 0
	}
	if y1 > yCenter {
		dy = y1 - yCenter
	} else if yCenter > y2 {
		dy = yCenter - y2
	} else {
		dy = 0
	}
	return dx*dx+dy*dy <= radius*radius
}

func (node *TreeNode) findGameObject(radius int, xCenter int, yCenter int) *LinkList {
	result := NewLinkList()
	if node.checkOverlap(radius, xCenter, yCenter) {
		node.objects.foreach(func(gameObject *GameObject) {
			x1 := xCenter
			x2 := gameObject.x
			y1 := yCenter
			y2 := gameObject.y
			if (x1-x2)*(x1-x2)+(y1-y2)*(y1-y2) < radius*radius {
				result.add(gameObject)
			}
		})

		if node.topLeft != nil {
			result.merge(node.topLeft.findGameObject(radius, xCenter, yCenter))
			result.merge(node.topRight.findGameObject(radius, xCenter, yCenter))
			result.merge(node.bottomLeft.findGameObject(radius, xCenter, yCenter))
			result.merge(node.bottomRight.findGameObject(radius, xCenter, yCenter))
		}
	}
	return result
}
