package entity

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type TrackID string

func NewTrackID() TrackID {
	letters := randomLetters(3)
	numbers := randomNumbers(3)

	return TrackID(fmt.Sprintf("%s-%s", letters, numbers))
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
	b := make([]rune, n)
	for i := range b {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
		if err != nil {
			panic(err)
		}

		b[i] = set[idx.Int64()]
	}
	return string(b)
}
