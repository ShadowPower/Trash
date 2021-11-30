package main

import "math"

type TreeNode struct {
	depth       int
	x           int
	y           int
	w           int
	h           int
	objects     *[]*GameObject
	topLeft     *TreeNode
	topRight    *TreeNode
	bottomLeft  *TreeNode
	bottomRight *TreeNode
	parent      *TreeNode
}

func NewTreeNode(depth int, x int, y int, w int, h int, objects *[]*GameObject, parent *TreeNode) *TreeNode {
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
	newGameObjects := append(*node.objects, gameObject)
	node.objects = &newGameObjects
}

func (node *TreeNode) removeGameObject(gameObject *GameObject) {
	newObjects := make([]*GameObject, len(*node.objects)-1)
	index := 0
	for _, origNodeObj := range *node.objects {
		if origNodeObj != gameObject {
			newObjects[index] = origNodeObj
			index += 1
		}
	}
	node.objects = &newObjects
}

func (node *TreeNode) update() {
	// 可以继续细分的条件
	if node.depth < 20 && node.w > 2 && node.h > 2 {
		if node.objects != nil && len(*node.objects) > 1 {
			var leftW = node.w / 2
			var topH = node.h / 2
			var rightX = node.x + leftW + 1
			var rightW = node.w - leftW - 1
			var bottomY = node.y + topH + 1
			var bottomH = node.h - topH - 1

			node.topLeft = NewTreeNode(node.depth+1, node.x, node.y, leftW, topH, &[]*GameObject{}, node)
			node.topRight = NewTreeNode(node.depth+1, rightX, node.y, rightW, topH, &[]*GameObject{}, node)
			node.bottomLeft = NewTreeNode(node.depth+1, node.x, bottomY, leftW, bottomH, &[]*GameObject{}, node)
			node.bottomRight = NewTreeNode(node.depth+1, rightX, bottomY, rightW, bottomH, &[]*GameObject{}, node)

			for _, gameObject := range *node.objects {
				if gameObject.x < rightX && gameObject.y < bottomY {
					gameObject.parent = node.topLeft
					var objs = append(*node.topLeft.objects, gameObject)
					node.topLeft.objects = &objs
				} else if gameObject.x >= rightX && gameObject.y < bottomY {
					gameObject.parent = node.topRight
					var objs = append(*node.topRight.objects, gameObject)
					node.topRight.objects = &objs
				} else if gameObject.x < rightX && gameObject.y >= bottomY {
					gameObject.parent = node.bottomLeft
					var objs = append(*node.bottomLeft.objects, gameObject)
					node.bottomLeft.objects = &objs
				} else if gameObject.x >= rightX && gameObject.y >= bottomY {
					gameObject.parent = node.bottomRight
					var objs = append(*node.bottomRight.objects, gameObject)
					node.bottomRight.objects = &objs
				}
			}

			node.objects = &[]*GameObject{}

			node.topLeft.update()
			node.topRight.update()
			node.bottomLeft.update()
			node.bottomRight.update()

			return
		}
	}
	// 更新不满足条件的 gameObject 的 parent
	for _, gameObject := range *node.objects {
		gameObject.parent = node
	}
}

func (node *TreeNode) optimize() int {
	if node.topLeft == nil && node.topRight == nil && node.bottomLeft == nil && node.bottomRight == nil {
		// 根节点，返回当前节点对象个数
		return len(*node.objects)
	} else {
		// 非根节点
		var childObjectNum = node.topLeft.optimize() + node.topRight.optimize() +
			node.bottomLeft.optimize() + node.bottomRight.optimize()
		if childObjectNum < 2 {
			if childObjectNum > 0 {
				// 子节点数据移到父节点
				var combineObjects = append(*node.objects, *node.topLeft.objects...)
				combineObjects = append(combineObjects, *node.topRight.objects...)
				combineObjects = append(combineObjects, *node.bottomLeft.objects...)
				combineObjects = append(combineObjects, *node.bottomRight.objects...)
				for _, gameObj := range combineObjects {
					gameObj.parent = node
				}
				node.objects = &combineObjects
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

func (node *TreeNode) findGameObject(radius int, xCenter int, yCenter int) []*GameObject {
	result := []*GameObject{}
	if node.checkOverlap(radius, xCenter, yCenter) {
		for _, gameObject := range *node.objects {
			x1 := float64(xCenter)
			x2 := float64(gameObject.x)
			y1 := float64(yCenter)
			y2 := float64(gameObject.y)
			distance := math.Sqrt(math.Pow(x1-x2, 2) + math.Pow(y1-y2, 2))
			if distance < float64(radius) {
				result = append(result, gameObject)
			}
		}

		if node.topLeft != nil {
			result = append(result, node.topLeft.findGameObject(radius, xCenter, yCenter)...)
			result = append(result, node.topRight.findGameObject(radius, xCenter, yCenter)...)
			result = append(result, node.bottomLeft.findGameObject(radius, xCenter, yCenter)...)
			result = append(result, node.bottomRight.findGameObject(radius, xCenter, yCenter)...)
		}
	}
	return result
}
