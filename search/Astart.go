package search

import "fmt"

const (
	width  = 24
	height = 32
)

type Node struct {
	X int
	Y int
	// F = G + H
	// F是总代价，G是从起点到当前点的代价，H是从当前点到终点的代价
	F, G, H int
	Parent  *Node
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func AStart(sneak, food [][]int) []int {
	openList := make([]*Node, 0)
	closeList := make([]*Node, 0)
	//将起始点设置优先级0,加入openList
	start := &Node{
		X: sneak[0][0],
		Y: sneak[0][1],
		G: 0,
		H: abs(sneak[0][0]-food[0][0]) + abs(sneak[0][1]-food[0][1]),
		F: abs(sneak[0][0]-food[0][0]) + abs(sneak[0][1]-food[0][1]),
	}
	end := &Node{
		X: food[0][0],
		Y: food[0][1],
	}
	openList = append(openList, start)

	for len(openList) > 0 {
		//找到openList中F值最小的点
		index := 0
		for i := 0; i < len(openList); i++ {
			if openList[i].F < openList[index].F {
				index = i
			}
		}
		cur := openList[index]

		//找到终点
		if cur.X == end.X && cur.Y == end.Y {
			path := make([]Node, 0)
			for cur.Parent != nil {
				path = append(path, *cur)
				cur = cur.Parent
			}
			return getPath(append(path, *cur))
		}
		//从openList中删除当前点
		openList = append(openList[:index], openList[index+1:]...)
		//将当前点加入closeList
		closeList = append(closeList, cur)
		//获取当前点的邻居
		neighborList := cur.getNeighbor(sneak[1:])
		for _, neighbor := range neighborList {
			if isExit(closeList, neighbor) {
				continue
			}
			neighbor.G = cur.G + 1
			neighbor.H = abs(neighbor.X-end.X) + abs(neighbor.Y-end.Y)
			neighbor.F = neighbor.G + neighbor.H
			//如果邻居不在openList中,则加入openList
			if !isExit(openList, neighbor) {
				neighbor.Parent = cur
				openList = append(openList, neighbor)
			}
		}
	}

	return nil
}

func isExit(list []*Node, node *Node) bool {
	for _, v := range list {
		if v.X == node.X && v.Y == node.Y {
			return true
		}
	}
	return false
}

func (n *Node) getNeighbor(mapData [][]int) []*Node {
	hash := make(map[string]bool)
	for _, v := range mapData {
		hash[fmt.Sprintf("%d-%d", v[0], v[1])] = true
	}
	list := make([]*Node, 0)
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			if n.X+i < 0 || n.X+i > width || n.Y+j < 0 || n.Y+j > height {
				continue
			}
			//如果是障碍物,即蛇身,则跳过
			if hash[fmt.Sprintf("%d-%d", n.X+i, n.Y+j)] {
				continue
			}

			if i*j == 0 {
				list = append(list, &Node{
					X: n.X + i,
					Y: n.Y + j,
				})
			}
		}
	}
	return list
}

func getPath(path []Node) []int {
	step := make([]int, 0)
	for i := len(path) - 1; i > 0; i-- {
		if path[i].X == path[i-1].X {
			if path[i].Y > path[i-1].Y {
				step = append(step, 2)
			} else {
				step = append(step, 3)
			}
		} else {
			if path[i].X > path[i-1].X {
				step = append(step, 0)
			} else {
				step = append(step, 1)
			}
		}
	}
	return step
}
