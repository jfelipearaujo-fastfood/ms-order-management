package entity

import (
	"fmt"
	"math/rand"
	"time"
)

type TrackID string

func NewTrackID() TrackID {
	letters := randomLetters(3)
	numbers := randomNumbers(3)

	return TrackID(fmt.Sprintf("%s%s", letters, numbers))
}

func NewTrackIDFrom(s string) TrackID {
	return TrackID(s)
}

var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numbers = []rune("0123456789")

func randomLetters(n int) string {
	return randomFromSet(letters, n)
}

func randomNumbers(n int) string {
	return randomFromSet(numbers, n)
}

func randomFromSet(set []rune, n int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)
	for i := range b {
		b[i] = set[rnd.Intn(len(set))]
	}
	return string(b)
}
