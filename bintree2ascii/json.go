package bintree2ascii

// Tree is a tree used for parsing json as binary tree.
//
// It implements our Interface, so we can easily generate ascii art for it.
type Tree struct {
	*TreeNode
}

type TreeNode struct {
	Name      string    `json:"name"`
	Left      *TreeNode `json:"left"`
	Right     *TreeNode `json:"right"`
	LeftEdge  string    `json:"leftEdge"`
	RightEdge string    `json:"rightEdge"`
}

func (t Tree) Left() Interface {
	if t.TreeNode.Left == nil {
		return nil
	}
	return Tree{t.TreeNode.Left}
}

func (t Tree) Right() Interface {
	if t.TreeNode.Right == nil {
		return nil
	}
	return Tree{t.TreeNode.Right}
}

func (t Tree) Key() string {
	return t.TreeNode.Name
}

func (t Tree) LeftEdge() string {
	return t.TreeNode.LeftEdge
}

func (t Tree) RightEdge() string {
	return t.TreeNode.RightEdge
}
