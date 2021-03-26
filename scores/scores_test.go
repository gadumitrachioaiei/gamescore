package scores

import (
	"fmt"
	"math/rand"
	"os/exec"
	"reflect"
	"sort"
	"testing"
	"time"
)

// TestScoresBST tests that the scrores form an BST tree,
// checking that inorder reverse produces the scores in descending order.
func TestScoresBST(t *testing.T) {
	s := New()
	_, sortedScores := generateScores(s)
	fmt.Println(s.root.ToAscii())
	fmt.Println(s.root.ToDot())
	assertBST(t, s, sortedScores)
}

// TestScoresRange tests Range scores.
func TestScoresRange(t *testing.T) {
	type testCase struct {
		scores   []Score
		position int
		count    int
		expected []Score
	}
	testCases := []testCase{
		{
			scores:   []Score{{0, 3}, {1, 1}, {2, 5}, {3, 2}},
			expected: []Score{{3, 2}, {1, 1}},
			position: 4,
			count:    1,
		},
		{
			scores:   []Score{{0, 7}, {9, 6}, {8, 6}, {6, 5}},
			expected: []Score{{9, 6}, {6, 5}},
			position: 4,
			count:    1,
		},
	}
	for _, tc := range testCases {
		s := New()
		for _, score := range tc.scores {
			s.Add(score)
		}
		if calculated := s.Range(tc.position, tc.count); !reflect.DeepEqual(calculated, tc.expected) {
			t.Fatalf("got count %d scores around position %d: \n%v\n expected: \n%v\n scores:\n%v\n",
				tc.count, tc.position, calculated, tc.expected, tc.scores)
		}
	}
}

// TestScoresRangeRandom tests Range scores using random test cases.
func TestScoresRangeRandom(t *testing.T) {
	s := New()
	scores, sortedScores := generateScores(s)
	count := 2
	for i := 0; i < len(sortedScores)+2; i++ {
		// calculate expected scores
		var expectedScores []Score
		top := i + count
		totalCount := 2*count + 1
		if top > len(scores) {
			totalCount = totalCount - (top - len(sortedScores))
			top = len(sortedScores)
		}
		if totalCount > top {
			totalCount = top
		}
		expectedScores = sortedScores[:top][top-totalCount:]
		if len(expectedScores) == 0 {
			expectedScores = nil
		}
		if calculated := s.Range(i, count); !reflect.DeepEqual(calculated, expectedScores) {
			t.Fatalf("got count %d scores around position %d: \n%v\n expected: \n%v\n scores:\n%v\n",
				count, i, calculated, expectedScores, scores)
		}
	}
}

// TestScoresTopRandom tests Top scores using random test cases.
func TestScoresTopRandom(t *testing.T) {
	s := New()
	scores, sortedScores := generateScores(s)
	for i := 0; i < len(sortedScores)+2; i++ {
		var expectedScores []Score
		if i >= len(sortedScores) {
			expectedScores = sortedScores
		} else if i > 0 {
			expectedScores = sortedScores[:i]
		}
		if calculated := s.Top(i); !reflect.DeepEqual(calculated, expectedScores) {
			t.Fatalf("got top %d scores: \n%v\n expected: \n%v\n, scores: %v", i, calculated, expectedScores, scores)
		}
	}
}

func TestUpdateRoot(t *testing.T) {
	s := New()
	s.Add(Score{User: 1, Value: 2})
	s.Add(Score{User: 2, Value: 1})
	s.Add(Score{User: 3, Value: 10})
	s.Add(Score{User: 4, Value: 8})
	s.Add(Score{User: 5, Value: 9})
	s.Add(Score{User: 6, Value: 6})
	s.Add(Score{User: 7, Value: 7})
	//fmt.Println(s.root.ToDot())
	text := s.root.ToDot()
	fmt.Println(text)
	comm := fmt.Sprintf(`echo '%s' | graph-easy --as_box`, text)
	fmt.Println(comm)
	result, err := exec.Command("bash", "-c", comm).CombinedOutput()
	if err != nil {
		t.Fatal(err, string(result))
	}
	fmt.Println(string(result))
	fmt.Println(s.root.ToAscii())
	//spew.Config.Indent = "\t"
	//spew.Dump(s.root)
}

// assertTree asserts that a tree
func assertTreeUpdate(t *testing.T, s *Scores, score Score) {
	sortedScores := inOrder(s.root)
	sortedScores = updateSortedScores(sortedScores, score)
	newScore, err := s.Update(score)
	if err != nil {
		t.Fatalf("updating existing user: %v", err)
	}
	expected := Score{User: 1, Value: 2}
	if newScore != expected {
		t.Fatalf("got new score: %v, expected: %v", newScore, expected)
	}
	assertBST(t, s, sortedScores)
	if node, ok := s.users[score.User]; !(ok && node.score == newScore.Value && node.user == newScore.User) {
		t.Fatalf("got mapped user score: %v, expected: %v", Score{User: node.user, Value: node.score}, newScore)
	}
}

func updateSortedScores(scores []Score, newScore Score) []Score {
	for i := 0; i < len(scores); i++ {
		if scores[i].User == newScore.User {
			newScore.Value += scores[i].Value
			scores[i] = scores[len(scores)-1]
			scores = scores[:len(scores)-1]
			break
		}
	}
	scores = append(scores, newScore)
	return sortScores(scores)
}

// generateScores attaches 10 random scores to the tree and returns them together with a sorted copy.
func generateScores(t *Scores) ([]Score, []Score) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	var scores []Score
	for i := 0; i < 10; i++ {
		score := Score{
			User:  i,
			Value: random.Intn(10),
		}
		scores = append(scores, score)
		t.Add(score)
	}
	return scores, sortScores(scores)
}

// sortScores returns sorted scores in descending order of scores and insertion time.
func sortScores(scores []Score) []Score {
	type ScoreWithIndex struct {
		Score
		Index int
	}
	scoresWithIndex := make([]ScoreWithIndex, len(scores))
	for i := 0; i < len(scores); i++ {
		scoresWithIndex[i] = ScoreWithIndex{scores[i], i}
	}
	sort.Slice(scoresWithIndex, func(i, j int) bool {
		if scoresWithIndex[i].Value == scoresWithIndex[j].Value {
			return scoresWithIndex[i].Index > scoresWithIndex[j].Index
		}
		return scoresWithIndex[i].Value > scoresWithIndex[j].Value
	})
	sortedScores := make([]Score, len(scoresWithIndex))
	for i := 0; i < len(scoresWithIndex); i++ {
		sortedScores[i] = scoresWithIndex[i].Score
	}
	return sortedScores
}

func assertBST(t *testing.T, s *Scores, sortedScores []Score) {
	scores := inOrder(s.root)
	if !reflect.DeepEqual(scores, sortedScores) {
		t.Fatalf("got scores:\n%v\n, expected:\n%v\n", scores, sortedScores)
	}
}

func inOrder(node *Node) []Score {
	if node == nil {
		return nil
	}
	scores := inOrder(node.right)
	scores = append(scores, Score{User: node.user, Value: node.score})
	scores = append(scores, inOrder(node.left)...)
	return scores
}
