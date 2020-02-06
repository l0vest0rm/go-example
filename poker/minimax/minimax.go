package minimax

import (
	"fmt"
	"math"
)

type Node struct {
	IsMiniNode bool //mini or max
	alpha      int
	beta       int
	Data       interface{}
	parent     *Node
	children   []*Node
}

// New 新跟节点
func New() *Node {
	node := Node{parent: nil, IsMiniNode: true, alpha: math.MinInt32, beta: math.MaxInt32}
	return &node
}

// AddChild Add a new node to structure, this node should have children and
// an unknown score
func (node *Node) AddChild(data interface{}) *Node {
	childNode := Node{parent: node, Data: data, alpha: node.alpha, beta: node.beta}
	childNode.IsMiniNode = !node.IsMiniNode
	if node.children == nil {
		node.children = make([]*Node, 0)
	}
	node.children = append(node.children, &childNode)
	return &childNode
}

// AddLeafChild 增加叶子节点
func (node *Node) AddLeafChild(data interface{}, score int) *Node {
	childNode := Node{parent: node, Data: data, alpha: score, beta: score}
	childNode.IsMiniNode = !node.IsMiniNode
	if node.children == nil {
		node.children = make([]*Node, 0)
	}
	node.children = append(node.children, &childNode)
	childNode.updateAlphaBeta()

	return &childNode
}

func (node *Node) NeedCut() bool {
	return node.beta <= node.alpha
}

// 更新父节点及以上节点的alpha和beta值
func (node *Node) updateAlphaBeta() {
	var score int
	curNode := node
	for curNode.parent != nil {
		if curNode.IsMiniNode {
			score = curNode.beta
		} else {
			score = curNode.alpha
		}

		curNode = curNode.parent
		if curNode.IsMiniNode {
			curNode.beta = min(score, curNode.beta)
		} else {
			curNode.alpha = max(score, curNode.alpha)
		}
	}
}

//是否叶子节点
func (node *Node) isLeaf() bool {
	return node.children == nil || len(node.children) == 0
}

// 获取最优选择的孩子节点的数据
func (node *Node) GetBestChildNode() *Node {
	var score int
	if node.IsMiniNode {
		score = node.beta
	} else {
		score = node.alpha
	}

	for _, childNode := range node.children {
		if childNode.IsMiniNode {
			if childNode.beta == score {
				return childNode
			}
		} else {
			if childNode.alpha == score {
				return childNode
			}
		}
	}

	return nil
}

// Print the node for debugging purposes
func (node *Node) Print(level int) {
	var padding = ""
	for j := 0; j < level; j++ {
		padding += " "
	}

	fmt.Println(padding, node.IsMiniNode, node.alpha, node.beta, node.Data)

	for _, cn := range node.children {
		level += 2
		cn.Print(level)
		level -= 2
	}
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
