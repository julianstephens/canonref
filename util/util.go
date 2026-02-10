package util

import (
	"fmt"
	"strconv"
)

const EnDash = "–"
const Hyphen = "-"

func Ptr[T any](v T) *T {
	return &v
}

func If[T any](cond bool, t, f T) T {
	if cond {
		return t
	}
	return f
}

type VerseRange struct {
	StartVerse int  `json:"start"`
	EndVerse   *int `json:"end,omitempty"`
}

func (v VerseRange) String() string {
	if v.EndVerse == nil {
		return strconv.Itoa(v.StartVerse)
	}
	return fmt.Sprintf("%d%s%d", v.StartVerse, EnDash, *v.EndVerse)
}
