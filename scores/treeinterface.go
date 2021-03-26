package scores

import (
	"strconv"

	"github.com/gadumitrachioaiei/gamescore/bintree2ascii"
)

func (s *Node) Left() bintree2ascii.Interface {
	if s.left == nil {
		return nil
	}
	return s.left
}

func (s *Node) Right() bintree2ascii.Interface {
	if s.right == nil {
		return nil
	}
	return s.right
}

func (s *Node) LeftEdge() string {
	return strconv.Itoa(s.lsize)
}

func (s *Node) RightEdge() string {
	return strconv.Itoa(s.rsize)
}

func (s *Node) Key() string {
	return strconv.Itoa(s.score) + " " + strconv.Itoa(s.user)
}
