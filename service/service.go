package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gadumitrachioaiei/gamescore/scores"
)

type Service struct {
	scores *scores.Scores
}

func New() *Service {
	return &Service{scores: scores.New()}
}

func (s *Service) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		s.AddScore(w, req)
		return
	}
	if strings.HasPrefix(req.URL.Path, "/scores/top") {
		s.Top(w, req)
		return
	}
	if strings.HasPrefix(req.URL.Path, "/scores/range") {
		s.Range(w, req)
		return
	}
}

type Score struct {
	User  int
	Total int
}

func (s *Service) AddScore(w http.ResponseWriter, req *http.Request) {
	var score Score
	if err := json.NewDecoder(req.Body).Decode(&score); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if score.User <= 0 {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	s.scores.Add(scores.Score{UserID: score.User, Value: score.Total})
	return
}

func (s *Service) Top(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	top, err := strconv.Atoi(req.Form.Get("top"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	topScores := s.scores.Top(int(top))
	if err := json.NewEncoder(w).Encode(topScores); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (s *Service) Range(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	position, err := strconv.Atoi(req.Form.Get("position"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	count, err := strconv.Atoi(req.Form.Get("count"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	topScores := s.scores.Range(position, count)
	if err := json.NewEncoder(w).Encode(topScores); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
