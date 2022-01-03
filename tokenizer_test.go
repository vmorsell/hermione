package main

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name   string
		r      io.Reader
		tokens []string
	}{
		{
			name:   "ok, one token",
			r:      strings.NewReader("Hello"),
			tokens: []string{"hello"},
		},
		{
			name:   "ok, multiple tokens",
			r:      strings.NewReader("Hello, world!"),
			tokens: []string{"hello", "world"},
		},
		{
			name:   "ok, multiple lines",
			r:      strings.NewReader("Hello, world!\nHow are you?"),
			tokens: []string{"hello", "world", "how", "are", "you"},
		},
		{
			name:   "ok, ending line break",
			r:      strings.NewReader("Hello, world!\nHow are you?\n"),
			tokens: []string{"hello", "world", "how", "are", "you"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewTokenizer(tt.r)
			for _, want := range tt.tokens {
				res, err := s.NextToken()
				require.Nil(t, err)
				require.Equal(t, want, res)
			}
		})
	}
}

func TestHasMoreTokens(t *testing.T) {
	tests := []struct {
		name string
		r    io.Reader
		res  bool
	}{
		{
			name: "no",
			r:    strings.NewReader(""),
			res:  false,
		},
		{
			name: "yes",
			r:    strings.NewReader("a"),
			res:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewTokenizer(tt.r)
			res := s.HasMoreTokens()
			require.Equal(t, tt.res, res)
		})
	}
}

// Tokenize the file tokenizer_test_corpus.txt and verify that the result
// corresponds with the tokens in tokenizer_test_tokens.txt
func TestTokenizeCorpus(t *testing.T) {
	corpusFile, err := os.Open("tokenizer_test_corpus.txt")
	require.Nil(t, err)
	defer corpusFile.Close()

	wantFile, err := os.Open("tokenizer_test_tokens.txt")
	require.Nil(t, err)
	defer wantFile.Close()

	want := make(map[string]int)
	scanner := bufio.NewScanner(wantFile)
	for scanner.Scan() {
		_, ok := want[scanner.Text()]
		if !ok {
			want[scanner.Text()] = 1
		} else {
			want[scanner.Text()] += 1
		}
	}

	s := NewTokenizer(corpusFile)

	got := make(map[string]int)
	for s.HasMoreTokens() {
		token, err := s.NextToken()
		require.Nil(t, err)
		_, ok := got[token]
		if !ok {
			got[token] = 1
		} else {
			got[token] += 1
		}
	}

	require.EqualValues(t, want, got)
}
