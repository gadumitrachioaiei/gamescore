package scores

import (
	"sync"
)

// Scores stores scores for our game and can answer queries about them.
//
// We store scores in a BST, with some additional metadata, so we can rank the scores.
//
// Thread safe.
type Scores struct {
	mu   sync.Mutex
	root *Node
}

// New returns a new Scores object
func New() *Scores {
	return &Scores{}
}

// Add adds a new score for the user in the s tree
func (s *Scores) Add(score Score) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.root == nil {
		s.root = &Node{
			score:  score.Value,
			userID: score.UserID,
		}
		return
	}
	s.root.Add(score)
}

// Top returns top scores in descending order.
func (s *Scores) Top(top int) []Score {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.root == nil {
		return nil
	}
	return s.root.Top(top)
}

// Range returns scores ranked between position-size and position+size, if they exist.
//
// The scores are sorted in descending order.
func (s *Scores) Range(position int, count int) []Score {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.root == nil {
		return nil
	}
	return s.root.Range(position, count)
}

// Node is a node for our scores tree.
type Node struct {
	score        int   // score of the user, used as key in our tree
	userID       int   // user that had the above score, used as value in our tree
	left, right  *Node // left and right children
	lsize, rsize int   // left and right subtree size
}

// Score represents a score, to be added or returned from our tree.
type Score struct {
	UserID int
	Value  int
}

// Add adds a new score for the user in the s tree.
func (s *Node) Add(score Score) {
	if score.Value < s.score {
		s.lsize++
		s.left.add(score, s)
		return
	}
	s.rsize++
	s.right.add(score, s)
}

func (s *Node) Update(score Score, rank int) int {
	return 0
}

// Top returns top scores, in descending order.
//
// If we have equal scores, the later ones are ranked higher.
func (s *Node) Top(top int) []Score {
	if top <= 0 {
		return nil
	}
	if top < s.rsize+1 {
		return s.right.Top(top)
	}
	var scores []Score
	if s.right != nil {
		s.right.inOrderReverse(&scores)
	}
	scores = append(scores, Score{s.userID, s.score})
	if top == s.rsize+1 || s.left == nil {
		return scores
	}
	top -= s.rsize + 1
	scores = append(scores, s.left.Top(top)...)
	return scores
}

// Range returns root ranked between position-size and position+size, if they exist.
//
// The root are sorted in descending order. If we have equal scores, the later ones are ranked higher.
func (s *Node) Range(position int, size int) []Score {
	var scores []Score
	s.search(1, position-size, position+size, &scores)
	return scores
}

// search searches scores ranked between startPos and endPos.
// startRank parameter represents the ranks underneath the current node, including the node itself
func (s *Node) search(startRank, startPos, endPos int, scores *[]Score) {
	// calculate the ranks of this node and its subtrees
	var (
		rightTreeRanks, leftTreeRanks [2]int
		nodeRank                      int
	)
	if s.rsize > 0 {
		rightTreeRanks[0] = startRank
		rightTreeRanks[1] = startRank + s.rsize - 1
	}
	nodeRank = startRank + s.rsize
	if s.lsize > 0 {
		leftTreeRanks[0] = nodeRank + 1
		leftTreeRanks[1] = nodeRank + s.lsize
	}
	// calculate where startPos and endPos fit, and walk the subtrees
	if s.rsize > 0 {
		r1, r2 := intersection(startPos, endPos, rightTreeRanks[0], rightTreeRanks[1])
		if r1 > 0 {
			s.right.search(rightTreeRanks[0], r1, r2, scores)
		}
	}
	if nodeRank >= startPos && nodeRank <= endPos {
		*scores = append(*scores, Score{
			UserID: s.userID,
			Value:  s.score,
		})
	}
	if s.lsize > 0 {
		r1, r2 := intersection(startPos, endPos, leftTreeRanks[0], leftTreeRanks[1])
		if r1 > 0 {
			s.left.search(leftTreeRanks[0], r1, r2, scores)
		}
	}
}

// inOrderReverse returns the users obtained by traversing the tree using in-order: right, parent, left.
func (s *Node) inOrderReverse(users *[]Score) {
	if s.right != nil {
		s.right.inOrderReverse(users)
	}
	*users = append(*users, Score{
		UserID: s.userID,
		Value:  s.score,
	})
	if s.left != nil {
		s.left.inOrderReverse(users)
	}
}

// add adds a new score for the user in the s tree, or attached directly to the parent
func (s *Node) add(user Score, parent *Node) {
	if s != nil {
		s.Add(user)
		return
	}
	newScore := &Node{
		score:  user.Value,
		userID: user.UserID,
	}
	if user.Value < parent.score {
		parent.left = newScore
		return
	}
	parent.right = newScore
}

// intersection returns the intersection between two intervals.
func intersection(s1, e1, s2, e2 int) (int, int) {
	if e1 < s2 || e2 < s1 {
		return -1, -1
	}
	if s2 >= s1 {
		if e2 <= e1 {
			return s2, e2
		}
		return s2, e1
	}
	if e1 <= e2 {
		return s1, e1
	}
	return s1, e2
}
