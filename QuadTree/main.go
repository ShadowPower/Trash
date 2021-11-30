package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	var size int
	fmt.Print("要生成的物体数量：")
	fmt.Scanln(&size)
	u := NewUniverse(4096, 4096)
	gameObjects := make([]*GameObject, size)
	for i := 0; i < size; i++ {
		gameObjects[i] = NewGameObject(rand.Intn(4096), rand.Intn(4096), i)
	}
	u.initGameObjects(gameObjects)

	// 调整位置
	testObject := gameObjects[size/2]
	testObject.x = 123
	testObject.y = 456
	u.moveGameObject(testObject)

	// 优化树结构
	u.elements.optimize()

	if size <= 10000 {
		for _, gameObject := range gameObjects {
			var obj = *gameObject
			print("对象 x:" + strconv.Itoa(obj.x) + ", y:" + strconv.Itoa(obj.y))
			print(" 位于区域 x:" + strconv.Itoa(obj.parent.x) + ", y:" + strconv.Itoa(obj.parent.y))
			println(" 内, 区域宽:" + strconv.Itoa(obj.parent.w) + ", 高:" + strconv.Itoa(obj.parent.h))
		}
	} else {
		println("数据太多，跳过输出")
	}

	var radius int
	print("查找半径：")
	fmt.Scanln(&radius)
	var xCenter int
	print("中心 x：")
	fmt.Scanln(&xCenter)
	var yCenter int
	print("中心 y：")
	fmt.Scanln(&yCenter)

	start := time.Now()
	findResult := u.findGameObject(radius, xCenter, yCenter)
	duration := time.Since(start)
	println("范围内的对象：")
	if len(findResult) <= 10000 {
		for _, result := range findResult {
			println(strconv.Itoa(result.x) + ", " + strconv.Itoa(result.y))
		}
	} else {
		println("数据太多，跳过输出")
	}

	print("查找耗时：")
	println(duration)

	var writeOut string
	print("保存树结构图像到 tree.png ？")
	fmt.Scanln(&writeOut)
	if strings.Contains(writeOut, "y") {
		writeOutImage(&u)
	}
}

func writeOutImage(universe *Universe) {
	resultImage := image.NewRGBA(image.Rect(0, 0, universe.width, universe.height))
	background := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	draw.Draw(resultImage, resultImage.Bounds(), &image.Uniform{C: background}, image.Point{}, draw.Src)

	layer := make([]*TreeNode, 0)
	layer = append(layer, universe.elements)
	for {
		nextLayer := make([]*TreeNode, 0)
		for _, node := range layer {
			if node.topLeft != nil {
				drawTreeNode(resultImage,0, node.topLeft)
				drawTreeNode(resultImage,1, node.topRight)
				drawTreeNode(resultImage,2, node.bottomLeft)
				drawTreeNode(resultImage,3, node.bottomRight)
				nextLayer = append(nextLayer, node.topLeft, node.topRight, node.bottomLeft, node.bottomRight)
			}
		}
		layer = nextLayer
		if len(layer) == 0 {
			break
		}
	}

	filename := "tree.png"
	pngFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	err = png.Encode(pngFile, resultImage)
	if err != nil {
		panic(err)
	}
}

func drawTreeNode(resultImage *image.RGBA, pos int, node *TreeNode)  {
	rect := image.Rect(node.x, node.y, node.x + node.w + 1, node.y + node.h + 1)
	var treeNodeColor color.RGBA
	colour := byte(node.depth * 12)
	switch pos {
	case 0:
		treeNodeColor = color.RGBA{R: colour, A: 255}
	case 1:
		treeNodeColor = color.RGBA{G: colour, A: 255}
	case 2:
		treeNodeColor = color.RGBA{B: colour, A: 255}
	case 3:
		treeNodeColor = color.RGBA{R: colour, G: colour, B: colour, A: 255}
	}

	draw.Draw(resultImage, rect, &image.Uniform{C: treeNodeColor}, image.Point{}, draw.Src)

	if len(*node.objects) > 0 {
		for _, gameObject := range *node.objects {
			resultImage.Set(gameObject.x, gameObject.y, color.RGBA{R: 255, B: 255, G: 255, A: 255})
		}
	}
}