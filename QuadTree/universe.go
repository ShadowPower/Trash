package main

type Universe struct {
	// 宽度
	width    int
	// 高度
	height 	 int
	elements *TreeNode
}

func NewUniverse(width int, height int) Universe {
	return Universe{
		width:    width,
		height:   height,
		elements: NewTreeNode(0, 0, 0, width, height, &[]*GameObject{}, nil),
	}
}

func (universe *Universe) initGameObjects(gameObjects []*GameObject) {
	universe.elements.objects = &gameObjects
	universe.elements.update()
}

func (universe *Universe) findTreeNode(x int, y int) *TreeNode {
	curNode := universe.elements
	if !(x > curNode.x && x < curNode.x+curNode.w && y > curNode.y && y < curNode.y+curNode.h) {
		panic("您飞出了世界")
	}

	for {
		// 取右下角的区域
		bottomRight := curNode.bottomRight
		// 如果没有子节点了
		if bottomRight == nil {
			return curNode
		}
		// 判断坐标位于哪个区域
		if x < bottomRight.x && y < bottomRight.y {
			curNode = curNode.topLeft
		} else if x >= bottomRight.x && y < bottomRight.y {
			curNode = curNode.topRight
		} else if x < bottomRight.x && y >= bottomRight.y {
			curNode = curNode.bottomLeft
		} else if x >= bottomRight.x && y >= bottomRight.y {
			curNode = bottomRight
		}
	}
}

func (universe *Universe) moveGameObject(gameObject *GameObject) {
	origNode := gameObject.parent
	if origNode != nil {
		// 不需要移动
		if gameObject.x > origNode.x && gameObject.x < origNode.x+origNode.w &&
			gameObject.y > origNode.y && gameObject.y < origNode.y+origNode.h {
			return
		}

		// 旧的树节点中删除当前元素
		origNode.removeGameObject(gameObject)
	}

	// 移动到新的节点下
	newNode := universe.findTreeNode(gameObject.x, gameObject.y)
	if origNode == newNode {
		return
	} else {
		newObjects := append(*newNode.objects, gameObject)
		newNode.objects = &newObjects
		newNode.update()
	}
}

func (universe *Universe) addGameObject(gameObject *GameObject) {
	treeNode := universe.findTreeNode(gameObject.x, gameObject.y)
	treeNode.addGameObject(gameObject)
	treeNode.update()
}

func (universe *Universe) removeGameObject(gameObject *GameObject) {
	treeNode := gameObject.parent
	if gameObject.parent != nil {
		treeNode.removeGameObject(gameObject)
	}
}

func (universe *Universe) findGameObject(radius int, x_center int, y_center int) []*GameObject {
	var result []*GameObject
	result = append(result, universe.elements.findGameObject(radius, x_center, y_center)...)
	return result
}
