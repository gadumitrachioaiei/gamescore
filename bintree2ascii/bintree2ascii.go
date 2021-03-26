// Package bintree2ascii generates ascii art for binary trees.
//
// If we will ever want to use fancy looks, look here:
// http://efanzh.org/tree-graph-generator/
// https://en.wikipedia.org/wiki/Box-drawing_character
package bintree2ascii

import (
	"encoding/json"
	"errors"
	"io"

	"gonum.org/v1/gonum/graph/formats/dot/ast"

	"gonum.org/v1/gonum/graph/formats/dot"
)

type Config struct {
	NodeWidth  int // width of a node's content
	NodeHeight int // height of a node's content
	EdgeHeight int // height of a edge
	Distance   int // Distance between to sibling nodes
	Sep        int // Distance between two consecutive nodes that are not siblings
}

type AsciiTree struct {
	Config
	//depth int // number of node levels

	levels []*Level
}

// NewAsciiTree returns a new ascii tree
func NewAsciiTree(config Config) *AsciiTree {
	return &AsciiTree{
		Config: config,
	}
}

type Interface interface {
	Left() Interface
	Right() Interface
	Key() string
	LeftEdge() string
	RightEdge() string
}

func (at *AsciiTree) FromInterface(i Interface) {
	treeLevels := levels(i)
	depth := len(treeLevels)
	asciiLevels := make([]*Level, 2*depth-1)
	asciiLevels[len(asciiLevels)-1] = LastLevel(LevelConfig{
		nodeWidth:  at.NodeWidth,
		nodeHeight: at.NodeHeight,
		distance:   at.Distance,
		sep:        at.Sep,
	}, at.asciiLevel(treeLevels[depth-1]))
	for i := len(asciiLevels) - 1; i >= 2; i -= 2 {
		// current tree level for this ascii level
		treeLevel := treeLevels[(i-2)/2]
		// node level
		asciiLevels[i-2] = ParentLevel(asciiLevels[i], at.asciiLevel(treeLevel))
		// edge level
		asciiLevels[i-1] = EdgeLevel(asciiLevels[i-2], asciiLevels[i], at.EdgeHeight, edgeLabels(treeLevel))
	}
	at.levels = asciiLevels
}

// FromJson generates an ascii tree from json, that needs to follow naming conventions:
// a node is described by: name, left, right, leftEdge, rightEdge.
func (at *AsciiTree) FromJson(r io.Reader) error {
	var tree Tree
	if err := json.NewDecoder(r).Decode(&tree); err != nil {
		return err
	}
	at.FromInterface(tree)
	return nil
}

// FromDot generates an ascii tree from adjacency data.
func (at *AsciiTree) FromDot(r io.Reader) error {
	f, err := dot.Parse(r)
	if err != nil {
		return err
	}
	if len(f.Graphs) != 1 {
		return errors.New("we need exactly one specified graph")
	}
	graph := f.Graphs[0]
	// we want to make the adjacency matrix
	type Edge struct {
		to    string
		label string
	}
	adj := make(map[string]*[2]Edge)
	for _, stmt := range graph.Stmts {
		edge, ok := stmt.(*ast.EdgeStmt)
		if !ok {
			continue
		}
		fromNode, ok := edge.From.(*ast.Node)
		if !ok {
			continue
		}
		toNode, ok := edge.To.Vertex.(*ast.Node)
		if !ok {
			continue
		}
		var edgeIndex int
		var label string
		for _, attr := range edge.Attrs {
			if attr.Key == "direction" && attr.Key == "right" {
				edgeIndex = 1
			} else if attr.Key == "label" {
				label = attr.Val
			}
		}
		adj[fromNode.ID][edgeIndex] = Edge{to: toNode.ID, label: label}
	}
	// find the root node
	var root string
	for _, v := range adj {
		var found bool
		for i := 0; i < len(v); i++ {
			if _, ok := adj[v[i].to]; !ok {
				root = v[i].label
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	// make up the tree
	var toTreeNode func(string) *TreeNode
	toTreeNode = func(node string) *TreeNode {
		if node == "" {
			return nil
		}
		return &TreeNode{
			Name:      node,
			Left:      toTreeNode(adj[node][0].to),
			Right:     toTreeNode(adj[node][1].to),
			LeftEdge:  adj[node][0].label,
			RightEdge: adj[node][1].label,
		}
	}
	at.FromInterface(Tree{toTreeNode(root)})
}

func (at *AsciiTree) Draw() []byte {
	var result []byte
	for i := 0; i < len(at.levels); i++ {
		result = append(result, at.levels[i].Draw()...)
	}
	return result
}

func (at *AsciiTree) asciiLevel(treeLevel []Interface) []Element {
	var nodes []Element
	for i := 0; i < len(treeLevel); i++ {
		var node *Node
		if treeLevel[i] == nil {
			node = NewInvisibleNode(at.NodeWidth, at.NodeHeight)
		} else {
			node = NewNode(treeLevel[i].Key(), at.NodeWidth, at.NodeHeight)
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func edgeLabels(level []Interface) []string {
	labels := make([]string, 2*len(level))
	for i := 0; i < len(level); i++ {
		if level[i] != nil {
			labels[2*i] = level[i].LeftEdge()
			labels[2*i+1] = level[i].RightEdge()
		}
	}
	return labels
}

func levels(i Interface) [][]Interface {
	var levels [][]Interface
	parentLevel := []Interface{i}
	levels = append(levels, parentLevel)
	for {
		isLastLevel := true
		var childLevel []Interface
		for i := 0; i < len(parentLevel); i++ {
			if parentLevel[i] == nil {
				childLevel = append(childLevel, nil, nil)
			} else {
				left, right := parentLevel[i].Left(), parentLevel[i].Right()
				if left != nil || right != nil {
					isLastLevel = false
				}
				childLevel = append(childLevel, left, right)
			}
		}
		if isLastLevel {
			break
		}
		levels = append(levels, childLevel)
		parentLevel = childLevel
	}
	return levels
}
