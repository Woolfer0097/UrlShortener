package shortener

import (
	"fmt"
	"testing"
)

func TestRandStringRunes_Length(t *testing.T) {
	for _, n := range []int{0, 1, 10, 100} {
		t.Run(fmt.Sprintf("Len%d", n), func(t *testing.T) {
			t.Parallel()
			s := RandStringRunes(n)
			if len(s) != n {
				t.Errorf("RandStringRunes(%d) length = %d", n, len(s))
			}
		})
	}
}

func TestRandStringRunes_Charset(t *testing.T) {
	t.Parallel()

	alphabetSet := make(map[rune]bool)
	for _, r := range alphabet {
		alphabetSet[r] = true
	}
	s := RandStringRunes(100)
	for _, r := range s {
		if !alphabetSet[r] {
			t.Errorf("RandStringRunes produced rune %q not in alphabet", r)
		}
	}
}

func TestGenerateShortUrl_Length(t *testing.T) {
	t.Parallel()
	out := GenerateShortUrl("https://github.com/Woolfer0097/")
	if len(out) != 10 {
		t.Errorf("GenerateShortUrl length = %d, want 10", len(out))
	}
}

func TestGenerateShortUrl_Charset(t *testing.T) {
	t.Parallel()
	alphabetSet := make(map[rune]bool)
	for _, r := range alphabet {
		alphabetSet[r] = true
	}
	out := GenerateShortUrl("https://github.com/Woolfer0097/")
	for _, r := range out {
		if !alphabetSet[r] {
			t.Errorf("GenerateShortUrl produced rune %q not in alphabet", r)
		}
	}
}
