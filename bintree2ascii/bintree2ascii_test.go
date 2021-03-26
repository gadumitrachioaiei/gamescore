package bintree2ascii

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

func makeTree() *TreeNode {
	n := TreeNode{
		Name:      "1",
		LeftEdge:  "3",
		RightEdge: "2",
	}
	n.Left = &TreeNode{
		Name:      "2",
		LeftEdge:  "1",
		RightEdge: "1",
		Left:      &TreeNode{Name: "4"},
		Right:     &TreeNode{Name: "5"},
	}
	n.Right = &TreeNode{
		Name:      "3",
		RightEdge: "1",
		Right:     &TreeNode{Name: "6"},
	}
	return &n
}

func TestNode_Draw(t *testing.T) {
	at := NewAsciiTree(Config{
		NodeWidth:  4,
		NodeHeight: 1,
		EdgeHeight: 3,
		Distance:   2,
		Sep:        1,
	})
	tree := makeTree()
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(tree); err != nil {
		t.Fatal(err)
	}
	r := io.TeeReader(&buf, os.Stdout)
	if err := at.FromJson(r); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", at.Draw())
	at.FromInterface(Tree{tree})
	fmt.Printf("%s\n", at.Draw())
}

func TestX(t *testing.T) {
	s := []byte("capsuna")
	suffix := []byte("1234242342342e213")
	fmt.Println(string(replaceSuffix(s, suffix)))
}
