package highscores

import ("slices"
        "sort")
type HighScores struct{
    scores []int
    latest int
    personalBest int
    topThree []int
}

// NewHighScores returns a new HighScores object.
func NewHighScores(scores []int) *HighScores {
	return &HighScores{
        scores: scores,
    }
}

// Scores returns all the scores.
func (s *HighScores) Scores() []int {
	return s.scores
}

// Latest returns the latest (last) score.
func (s *HighScores) Latest() int {
	last := s.scores[len(s.scores) - 1]
    return last
}

// PersonalBest returns the best (highest) score.
func (s *HighScores) PersonalBest() int {
	return slices.Max(s.scores)
}

// TopThree returns the top three scores.
func (s *HighScores) TopThree() []int {
	copi := append([]int{}, s.scores...)
    sort.Slice(copi, func(i, j int) bool {
        return copi[i] > copi[j]
    })
    if len(copi) <= 3 {
        return copi
    }
    return copi[:3]
}
