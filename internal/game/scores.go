package game

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ScoreEntry struct {
	Name  string
	Level int
	Score int
}

const scoresFile = "scores.txt"

func SaveScore(name string, level, score int) error {
	scores, _ := LoadScores()
	scores = append(scores, ScoreEntry{Name: name, Level: level, Score: score})

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	if len(scores) > 10 {
		scores = scores[:10]
	}

	f, err := os.Create(scoresFile)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, s := range scores {
		f.WriteString(s.Name + "|" + strconv.Itoa(s.Level) + "|" + strconv.Itoa(s.Score) + "\n")
	}

	return nil
}

func LoadScores() ([]ScoreEntry, error) {
	f, err := os.Open(scoresFile)
	if err != nil {
		return []ScoreEntry{}, nil
	}
	defer f.Close()

	var scores []ScoreEntry
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "|")
		if len(parts) == 3 {
			level, _ := strconv.Atoi(parts[1])
			score, _ := strconv.Atoi(parts[2])
			scores = append(scores, ScoreEntry{
				Name:  parts[0],
				Level: level,
				Score: score,
			})
		}
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	return scores, nil
}

func GetTopScores(n int) []ScoreEntry {
	scores, _ := LoadScores()
	if len(scores) > n {
		return scores[:n]
	}
	return scores
}
