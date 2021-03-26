package bintree2ascii

import (
	"bytes"
)

// Element is the interface for representing nodes and edges in a level.
type Element interface {
	Width() int           // The width as occupied by this element
	ContentWidth() int    // The width of content
	Next() ([]byte, bool) // Should return the nextLine line of the element and whether this is a valid line
	IsInvisible() bool    // whether this element is invisible or not.
}

// Node is a tree node.
type Node struct {
	key           string
	width, height int
	box           []byte
	isInvisible   bool

	offset int // current offset in the box
}

// NewNode returns a new node.
func NewNode(key string, width, height int) *Node {
	n := Node{key: key, width: width, height: height}
	n.fill()
	return &n
}

// NewInvisibleNode returns a node that will be invisible but still take space.
func NewInvisibleNode(width, height int) *Node {
	n := Node{width: width, height: height, isInvisible: true}
	n.fill()
	replaceWithSpaces(n.box)
	return &n
}

func (n *Node) IsInvisible() bool {
	return n.isInvisible
}

// Width returns the width as occupied by this node.
func (n *Node) Width() int {
	return n.width + 2
}

// ContentWidth returns the width of the content of the node.
func (n *Node) ContentWidth() int {
	return n.width
}

// Height returns the height as occupied by this node.
func (n *Node) Height() int {
	return n.height + 2
}

// Next returns the nextLine unread line of this node.
func (n *Node) Next() ([]byte, bool) {
	return nextLine(&n.offset, n.box)
}

func (n *Node) fill() {
	var border, content []byte
	defer func() {
		n.box = append(n.box, border...)
		n.box = append(n.box, '\n')
		n.box = append(n.box, content...)
		n.box = append(n.box, border...)
		n.box = append(n.box, '\n')
	}()
	// border
	border = make([]byte, 0, n.width+2)
	border = append(border, '+')
	border = append(border, bytes.Repeat([]byte{'-'}, n.width)...)
	border = append(border, '+')
	// content
	for i := 0; i < n.height; i++ {
		content = append(content, '|')
		content = append(content, bytes.Repeat([]byte{' '}, n.width)...)
		content = append(content, '|')
		content = append(content, '\n')
	}
	// put the Name in our box, in case it is not empty
	var start int
	for i := 0; i < len(content); i++ {
		if start >= len(n.key) {
			break
		}
		if content[i] != ' ' {
			continue
		}
		content[i] = n.key[start]
		start++
	}
}

type Edge struct {
	contentWidth, contentHeight int // content width and height
	box                         []byte
	isInvisible                 bool
	label                       string

	offset int // current offset in the box
}

func NewLeftEdge(width, height int, label string, isInvisible bool) *Edge {
	e := Edge{contentWidth: width - 2, contentHeight: height, label: label, isInvisible: isInvisible}
	e.drawLeft()
	if isInvisible {
		replaceWithSpaces(e.box)
	}
	return &e
}

func NewRightEdge(width, height int, label string, isInvisible bool) *Edge {
	e := Edge{contentWidth: width - 2, contentHeight: height, label: label, isInvisible: isInvisible}
	e.drawRight()
	if isInvisible {
		replaceWithSpaces(e.box)
	}
	return &e
}

func (e *Edge) IsInvisible() bool {
	return e.isInvisible
}

func (e *Edge) Width() int {
	return e.contentWidth + 2
}

func (e *Edge) ContentWidth() int {
	return e.contentWidth
}

func (e *Edge) Next() ([]byte, bool) {
	return nextLine(&e.offset, e.box)
}

func (e *Edge) drawLeft() {
	var horizontal, parentAnchor, childAnchor []byte
	defer func() {
		var buf []byte
		buf = append(buf, horizontal...)
		buf = append(buf, parentAnchor...)
		buf = append(buf, '\n')
		buf = append(buf, childAnchor...)
		e.box = buf
	}()
	horizontal = append(horizontal, ' ')
	for i := 0; i < e.contentWidth; i++ {
		horizontal = append(horizontal, '_')
	}
	parentAnchor = []byte{'|'}
	// build child anchor
	childAnchorLine := []byte{'|'}
	childAnchorLine = append(childAnchorLine, bytes.Repeat([]byte{' '}, e.contentWidth+1)...)
	for i := 0; i < e.contentHeight-1; i++ {
		if i == 0 {
			childAnchor = append(childAnchor, replacePrefix(childAnchorLine, []byte(e.label))...)
		} else {
			childAnchor = append(childAnchor, childAnchorLine...)
		}
		childAnchor = append(childAnchor, '\n')
	}
}

func (e *Edge) drawRight() {
	var horizontal, parentAnchor, childAnchor []byte
	defer func() {
		var buf []byte
		buf = append(buf, parentAnchor...)
		buf = append(buf, horizontal...)
		buf = append(buf, '\n')
		buf = append(buf, childAnchor...)
		e.box = buf
	}()
	parentAnchor = []byte{'|'}
	for i := 0; i < e.contentWidth; i++ {
		horizontal = append(horizontal, '_')
	}
	horizontal = append(horizontal, ' ')
	// build child anchor
	childAnchorLine := bytes.Repeat([]byte{' '}, e.contentWidth+1)
	childAnchorLine = append(childAnchorLine, '|')
	for i := 0; i < e.contentHeight-1; i++ {
		if i == 0 {
			childAnchor = append(childAnchor, replaceSuffix(childAnchorLine, []byte(e.label))...)
		} else {
			childAnchor = append(childAnchor, childAnchorLine...)
		}
		childAnchor = append(childAnchor, '\n')
	}
}

// replaceWithSpaces replaces all bytes except new line with space.
func replaceWithSpaces(data []byte) {
	for i := 0; i < len(data); i++ {
		if data[i] != '\n' {
			data[i] = ' '
		}
	}
}

// replacePrefix returns a copy of s, with the first chars from s replaced with prefix.
func replacePrefix(s, prefix []byte) []byte {
	b := append(s[:0:0], s...)
	copy(b, prefix)
	return b
}

// replaceSuffix returns a copy of s, with the last chars from s replaced with suffix.
func replaceSuffix(s, suffix []byte) []byte {
	if len(suffix) > len(s) {
		return append(suffix[:0:0], suffix[:len(s)]...)
	}
	b := append(s[:0:0], s...)
	copy(b[len(b)-len(suffix):], suffix)
	return b
}

// nextLine returns the next line from the box starting from the offset.
// it also updates offset to be the next indexed to start search from.
func nextLine(offset *int, box []byte) ([]byte, bool) {
	box = box[*offset:]
	off := bytes.IndexByte(box, '\n')
	if off == -1 {
		return nil, false
	}
	*offset += off + 1
	return box[:off], true
}
