package gomcts

import (
	"fmt"
	"math"
)

type miniMaxNode struct {
	IsMiniNode    bool //mini or max
	alpha         float64
	beta          float64
	causingAction Action
	state         GameState
	parent        *miniMaxNode
	children      []*miniMaxNode
}

func MiniMaxSearch(state GameState) Action {
	root := newMiniMaxNode(state)
	expandNode(root)

	childNode := root.GetBestChildNode()
	return childNode.causingAction
}

func expandNode(node *miniMaxNode) {
	actions := node.state.GetLegalActions()
	for i := 0; i < len(actions); i++ {
		if node.NeedCut() {
			return
		}

		newState := actions[i].ApplyTo(node.state)
		if newState.IsGameEnded() {
			score, _ := newState.EvaluateGame()
			node.AddLeafChild(newState, actions[i], float64(score))
		} else {
			expandNode(node.AddChild(newState, actions[i]))
		}
	}
}

// New 新跟节点
func newMiniMaxNode(state GameState) *miniMaxNode {
	node := miniMaxNode{parent: nil, IsMiniNode: true, alpha: math.MinInt32, beta: math.MaxInt32, state: state}
	return &node
}

// AddChild Add a new node to structure, this node should have children and
// an unknown score
func (node *miniMaxNode) AddChild(state GameState, action Action) *miniMaxNode {
	childNode := miniMaxNode{parent: node, state: state, causingAction: action, alpha: node.alpha, beta: node.beta}
	childNode.IsMiniNode = !node.IsMiniNode
	if node.children == nil {
		node.children = make([]*miniMaxNode, 0)
	}
	node.children = append(node.children, &childNode)
	return &childNode
}

// AddLeafChild 增加叶子节点
func (node *miniMaxNode) AddLeafChild(state GameState, action Action, score float64) *miniMaxNode {
	childNode := miniMaxNode{parent: node, state: state, causingAction: action, alpha: score, beta: score}
	childNode.IsMiniNode = !node.IsMiniNode
	if node.children == nil {
		node.children = make([]*miniMaxNode, 0)
	}
	node.children = append(node.children, &childNode)
	childNode.updateAlphaBeta()

	return &childNode
}

func (node *miniMaxNode) NeedCut() bool {
	return node.beta <= node.alpha
}

// 更新父节点及以上节点的alpha和beta值
func (node *miniMaxNode) updateAlphaBeta() {
	var score float64
	curNode := node
	for curNode.parent != nil {
		if curNode.IsMiniNode {
			score = curNode.beta
		} else {
			score = curNode.alpha
		}

		curNode = curNode.parent
		if curNode.IsMiniNode {
			curNode.beta = math.Min(score, curNode.beta)
		} else {
			curNode.alpha = math.Max(score, curNode.alpha)
		}
	}
}

//是否叶子节点
func (node *miniMaxNode) isLeaf() bool {
	return node.children == nil || len(node.children) == 0
}

// 获取最优选择的孩子节点的数据
func (node *miniMaxNode) GetBestChildNode() *miniMaxNode {
	var score float64
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
func (node *miniMaxNode) Print(level int) {
	var padding = ""
	for j := 0; j < level; j++ {
		padding += " "
	}

	fmt.Println(padding, node.IsMiniNode, node.alpha, node.beta, node.causingAction)

	for _, cn := range node.children {
		level += 2
		cn.Print(level)
		level -= 2
	}
}
