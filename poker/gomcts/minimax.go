package gomcts

import (
	"fmt"
	"math"
)

type miniMaxNode struct {
	depth         int
	IsMiniNode    bool //mini or max
	alpha         float64
	beta          float64
	causingAction Action
	state         GameState
	parent        *miniMaxNode
	children      []*miniMaxNode
}

func MiniMaxSearch(state GameState) (Action, bool) {
	root := newMiniMaxNode(state)
	bestChild := expandNode(root)
	//root.Print()

	if bestChild.beta > 0 && bestChild.beta < 2 {
		return bestChild.causingAction, true
	} else {
		return bestChild.causingAction, false
	}
}

func spPrint(node *miniMaxNode) {
	child := node.children[1]
	child.Print()
}

// return best choice child
func expandNode(node *miniMaxNode) *miniMaxNode {
	var childNode *miniMaxNode
	var bestChild *miniMaxNode
	actions := node.state.GetLegalActions()
	for i := 0; i < len(actions); i++ {
		newState := actions[i].ApplyTo(node.state)
		if newState.IsGameEnded() {
			score, _ := newState.EvaluateGame()
			childNode = node.AddLeafChild(newState, actions[i], float64(score))
		} else {
			childNode = node.AddChild(newState, actions[i])
			expandNode(childNode)
		}

		if node.IsMiniNode {
			if childNode.alpha < node.beta {
				node.beta = childNode.alpha
				bestChild = childNode
			}

			node.beta = math.Min(childNode.alpha, node.beta)
		} else {
			if childNode.beta > node.alpha {
				node.alpha = childNode.beta
				bestChild = childNode
			}
		}

		if node.NeedCut() {
			break
		}
	}

	return bestChild
}

// New 新跟节点
func newMiniMaxNode(state GameState) *miniMaxNode {
	node := miniMaxNode{depth: 0, parent: nil, IsMiniNode: false, alpha: math.MinInt32, beta: math.MaxInt32, state: state}
	return &node
}

// AddChild Add a new node to structure, this node should have children and
// an unknown score
func (node *miniMaxNode) AddChild(state GameState, action Action) *miniMaxNode {
	childNode := miniMaxNode{depth: node.depth + 1, parent: node, state: state, causingAction: action, alpha: node.alpha, beta: node.beta}
	childNode.IsMiniNode = !node.IsMiniNode
	/*if childNode.depth == 10 || childNode.depth == 11 {
		fmt.Printf("\nAddChild,%d,%v,%p,%p,%f,%f,%v", childNode.depth, childNode.IsMiniNode, &childNode, &childNode.parent, childNode.alpha, childNode.beta, childNode.causingAction)
	}*/
	if node.children == nil {
		node.children = make([]*miniMaxNode, 0)
	}
	node.children = append(node.children, &childNode)
	return &childNode
}

// AddLeafChild 增加叶子节点
func (node *miniMaxNode) AddLeafChild(state GameState, action Action, score float64) *miniMaxNode {
	childNode := miniMaxNode{depth: node.depth + 1, parent: node, state: state, causingAction: action, alpha: score, beta: score}
	childNode.IsMiniNode = !node.IsMiniNode

	if node.children == nil {
		node.children = make([]*miniMaxNode, 0)
	}
	node.children = append(node.children, &childNode)
	//fmt.Printf("\nAddLeafChild,%d,%v,%p,%p,%f,%f,%v", childNode.depth, childNode.IsMiniNode, &childNode, &childNode.parent, childNode.alpha, childNode.beta, childNode.causingAction)
	//childNode.updateAlphaBeta()

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
		/*if curNode.depth == 10 || curNode.depth == 11 {
			fmt.Printf("\nbefore,%d,%v,%p,%p,%f,%f,%v", curNode.depth, curNode.IsMiniNode, &curNode, &curNode.parent, curNode.alpha, curNode.beta, curNode.causingAction)
		}*/

		if curNode.IsMiniNode {
			curNode.beta = math.Min(score, curNode.beta)
		} else {
			curNode.alpha = math.Max(score, curNode.alpha)
		}

		/*if curNode.depth == 10 || curNode.depth == 11 {
			fmt.Printf("\nafter,%d,%v,%p,%p,%f,%f,%v", curNode.depth, curNode.IsMiniNode, &curNode, &curNode.parent, curNode.alpha, curNode.beta, curNode.causingAction)
		}*/
	}
}

//是否叶子节点
func (node *miniMaxNode) isLeaf() bool {
	return node.children == nil || len(node.children) == 0
}

// Print the node for debugging purposes
func (node *miniMaxNode) Print() {
	var padding = ""
	for j := 0; j < node.depth; j++ {
		padding += " "
	}

	fmt.Println(padding, node.depth, node.IsMiniNode, node.alpha, node.beta, node.causingAction)

	for _, cn := range node.children {
		cn.Print()
	}
}
