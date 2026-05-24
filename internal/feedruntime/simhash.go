package feedruntime

import (
	"hash/fnv"
	"math/bits"
	"regexp"
	"strings"
	"unicode"
)

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

// normalizeSimhashText strips HTML tags, lowercases, and collapses whitespace.
// Produces a clean token stream for 2-gram SimHash computation.
func normalizeSimhashText(s string) string {
	s = htmlTagRe.ReplaceAllString(s, " ")
	s = strings.ToLower(s)
	var b strings.Builder
	b.Grow(len(s))
	prevSpace := true // start as true to trim leading space
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !prevSpace {
				b.WriteByte(' ')
				prevSpace = true
			}
		} else {
			b.WriteRune(r)
			prevSpace = false
		}
	}
	result := b.String()
	if len(result) > 0 && result[len(result)-1] == ' ' {
		result = result[:len(result)-1]
	}
	return result
}

// computeSimhash returns a 64-bit SimHash fingerprint of the input text.
// It uses 2-gram (bigram) character features and FNV-64a hashing.
//
// Algorithm outline:
//  1. Split text into overlapping character 2-grams
//  2. For each gram, compute its FNV-64a hash h
//  3. For each bit position i: if bit i of h is 1, increment v[i]; else decrement v[i]
//  4. Build fingerprint: set bit i = 1 if v[i] > 0
func computeSimhash(text string) uint64 {
	runes := []rune(text)
	if len(runes) == 0 {
		return 0
	}

	var v [64]int

	h := fnv.New64a()
	for i := 0; i+1 < len(runes); i++ {
		h.Reset()
		gram := string(runes[i : i+2])
		_, _ = h.Write([]byte(gram))
		hash := h.Sum64()
		for bit := 0; bit < 64; bit++ {
			if hash&(1<<uint(bit)) != 0 {
				v[bit]++
			} else {
				v[bit]--
			}
		}
	}

	// Handle single-rune text: use the rune itself as a 1-gram
	if len(runes) == 1 {
		h.Reset()
		_, _ = h.Write([]byte(string(runes[0])))
		hash := h.Sum64()
		for bit := 0; bit < 64; bit++ {
			if hash&(1<<uint(bit)) != 0 {
				v[bit]++
			} else {
				v[bit]--
			}
		}
	}

	var fingerprint uint64
	for bit := 0; bit < 64; bit++ {
		if v[bit] > 0 {
			fingerprint |= 1 << uint(bit)
		}
	}
	return fingerprint
}

// hammingDistance counts the number of bit positions where a and b differ.
func hammingDistance(a, b uint64) int {
	return bits.OnesCount64(a ^ b)
}
