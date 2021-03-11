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
		s.replaceLeaf()
		return nil
	} else if s.left == nil {
		return s.replaceWith(s.right)
	} else if s.right == nil {
		return s.replaceWith(s.left)
	}
	return s.replaceWithDescendant()
}

// replaceLeaf replaces this leaf node.
func (s *Node) replaceLeaf() {
	defer s.parent.decrementSize()
	if s.parent.left != nil && s.parent.left.user == s.user {
		s.parent.left = nil
		s.parent.lsize = 0
		return
	}
	s.parent.right = nil
	s.parent.rsize = 0
}

// replaceWithRight replaces this node with the given node.
// returns the new node.
func (s *Node) replaceWith(n *Node) *Node {
	if s.parent.left != nil && s.parent.left.user == s.user {
		s.parent.left = n
	} else {
		s.parent.right = n
	}
	n.decrementSize()
	return n
}

// right replaces this node from the tree with a new one such that BST property is preserved.
//
// It returns the new node.
//
// The node must have both children.
func (s *Node) replaceWithDescendant() *Node {
	defer s.nullify()
	replacer := s.right.walkLeft()
	if s.right.left == nil {
		replacer = s.right
		s.right.left = s.left
		s.right.lsize = s.lsize
		s.decrementSize()
	} else {
		// its right tree becomes replacer of the parent
		replacer.parent.left = replacer.right
		replacer.parent.lsize -= 1
		replacer.parent.decrementSize()
		// its right tree is set to be node's right
		replacer.right = s.right
		// its replacer tree is set to be node's replacer
		replacer.left = s.left
		replacer.lsize = s.lsize
		replacer.rsize = s.rsize
	}
	if s.parent != nil {
		if s.parent.left != nil && s.parent.left.user == s.user {
			s.parent.left = replacer
		} else {
			s.parent.right = replacer
		}
	}
	return replacer
}

// decrementSize decrements size of ancestors of this node.
func (s *Node) decrementSize() {
	if s.parent == nil {
		return
	}
	// identify if the parent is to the left or right
	if s.parent.left != nil && s.parent.left.user == s.user {
		s.parent.lsize -= 1
	} else {
		s.parent.rsize -= 1
	}
	s.parent.decrementSize()
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
