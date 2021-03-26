package scores

import (
	"fmt"
)

// Update updates score for existing user and returns new score.
//
// Returns new score.
func (s *Scores) Update(score Score) (Score, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[score.User]; !ok {
		return Score{}, fmt.Errorf("user cannot be found: %d", score.User)
	}
	if score.User == s.root.user {
		return s.updateRoot(score)
	}
	node := s.users[score.User]
	node.replace()
	score.Value += node.score
	node = s.root.Add(score)
	s.users[score.User] = node
	return score, nil
}

// updateRoot updates root node with new score.
// Returns new score.
func (s *Scores) updateRoot(score Score) (Score, error) {
	score.Value += s.root.score
	if s.root.left == nil && s.root.right == nil {
		s.root.score = score.Value
		return score, nil
	}
	oldRoot := s.root
	defer oldRoot.nullify()
	if s.root.right == nil {
		s.root = s.root.left
	} else if s.root.left == nil {
		s.root = s.root.right
	} else {
		s.root = s.root.replaceWithDescendant()
	}
	delete(s.users, score.User)
	s.users[score.User] = s.root.Add(score)
	return score, nil
}

// replace replaces this node with another one, such that BST property is preserved.
//
// returns the replacement, which is nil if s is a leaf node.
func (s *Node) replace() *Node {
	defer s.nullify()
	if s.left == nil && s.right == nil {
		replaceParentConnection(s, nil)
		return nil
	}
	if s.left == nil {
		replaceParentConnection(s, s.right)
		return s.right
	}
	if s.right == nil {
		replaceParentConnection(s, s.left)
		return s.left
	}
	return s.replaceWithDescendant()
}

// right replaces this node from the tree with a new one such that BST property is preserved.
//
// It returns the new node.
//
// The node must have both children.
func (s *Node) replaceWithDescendant() *Node {
	defer s.nullify()
	if s.right.left == nil {
		replacer := s.right
		replacer.left = s.left
		replacer.lsize = s.lsize
		if s.parent != nil {
			replaceParentConnection(s, replacer)
		}
		return replacer
	}
	replacer := s.right.walkLeft()
	replaceParentConnection(replacer, replacer.right)
	replaceNode(s, replacer)
	return replacer
}

// subtractSize subtracts size from ancestors of this node.
func (s *Node) subtractSize(size int) {
	if size == 0 || s.parent == nil {
		return
	}
	// identify if the parent is to the left or right
	if s.parent.left != nil && s.parent.left.user == s.user {
		s.parent.lsize -= size
	} else {
		s.parent.rsize -= size
	}
	s.parent.subtractSize(size)
}

// walkLeft walks to the left of this node all the way and returns the last node.
func (s *Node) walkLeft() *Node {
	if s.left == nil {
		return s
	}
	return s.left.walkLeft()
}

// nullify removes all pointers of this node.
//
// You may need to call this method after you remove the node from the tree.
func (s *Node) nullify() {
	s.left, s.right, s.parent = nil, nil, nil
}

// replaceParentConnection replaces the connection between s's parent to s with s's parent to newChild.
// we assume that s has a parent.
func replaceParentConnection(s *Node, newChild *Node) {
	var newSize int
	if newChild != nil {
		newSize = newChild.lsize + newChild.rsize + 1
	}
	subtractedSize := s.parent.lsize + s.parent.rsize
	if s.parent.left != nil && s.parent.left.user == s.user {
		s.parent.left = newChild
		s.parent.lsize = newSize
		return
	}
	s.parent.right = newChild
	s.parent.rsize = newSize
	s.parent.subtractSize(subtractedSize - (s.parent.lsize + s.parent.rsize))
}

// replaceNode replaces node s with node n.
// you need to make sure after this call,
func replaceNode(s *Node, n *Node) {
	n.right = s.right
	n.rsize = s.rsize
	n.left = s.left
	n.lsize = s.lsize
	if s.parent != nil {
		replaceParentConnection(s, n)
	}
}
