package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
)

var stopBytes = []byte{' ', '\n'}

var patterns = []string{
	`https?:\/\/[a-z0-9\/\?&\.,=\-_:#\+%@!]+`, // URL:s
	`[a-z0-9\.\-_]+@[a-z0-9\.\-]+\.[a-z]+`,    // E-mail addresses.
	`[a-z]+\.[a-z]+\.`,                        // Abbreviations.
	`[a-z0-9_]+\.[a-z]+`,                      // File names.
	`([a-z0-9]+\-)+[a-z0-9]+`,                 // Words and numbers with dashes in them.
	`([0-9]+\.)+[0-9]+`,                       // Numbers with dots.
	`([0-9]+\-)+[0-9]+`,                       // Numbers with dashes.
	`([0-9]+,)+[0-9]+`,                        // Numbers with commas.
	`[0-9]+\/[0-9]+`,                          // 24/7 etc.
	`[a-z]+'[a-z]+`,                           // he's etc.
	`'em`,                                     // 'em.
	`([0-9]+:)+[0-9]+`,                        // Digital time representations.
	`\$[0-9]+`,                                // USD amounts.
	`\+[0-9]+`,                                // Phone numbers
	`[a-z0-9]+`,                               // Default token.
}

type Tokenizer interface {
	HasMoreWords() bool
	NextWord() ([]byte, error)
	HasMoreTokens() bool
	NextToken() (string, error)
}

type tokenizer struct {
	r        *bufio.Reader
	queue    [][]byte
	patterns []regexp.Regexp
}

func NewTokenizer(reader io.Reader) Tokenizer {
	return &tokenizer{
		r:        bufio.NewReader(reader),
		patterns: tokenPatterns(),
	}
}

// HasMoreWords returns if the reader had unread words.
func (t *tokenizer) HasMoreWords() bool {
	if _, err := t.r.Peek(1); err != nil {
		if err == io.EOF {
			return false
		}
	}
	return true
}

// NextWord returns the next available word from the reader.
func (t *tokenizer) NextWord() ([]byte, error) {
	var out bytes.Buffer

readByte:
	for {
		b, err := t.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("read byte: %w", err)
		}
		for _, s := range stopBytes {
			if s == b {
				break readByte
			}
		}
		c := bytes.ToLower([]byte{b})[0]
		out.WriteByte(c)
	}

	if out.Len() == 0 {
		if !t.HasMoreWords() {
			return nil, nil
		}
		out.Reset()
		goto readByte
	}

	return out.Bytes(), nil
}

func (t *tokenizer) HasMoreTokens() bool {
	if len(t.queue) != 0 {
		return true
	}
	if t.HasMoreWords() {
		return true
	}
	return false
}

// NextToken returns the next token.
func (t *tokenizer) NextToken() (string, error) {
	for t.HasMoreTokens() {
		// Make sure we have words in the queue.
		for len(t.queue) == 0 && t.HasMoreWords() {
			word, err := t.NextWord()
			if err != nil {
				return "", fmt.Errorf("read next word: %w", err)
			}
			t.queue = append(t.queue, word)
		}

		// Return if the queue is empty.
		if len(t.queue) == 0 {
			return "", nil
		}

		word := t.queue[0]
		t.queue = t.queue[1:]

		token, err := t.TokenFromWord(word)
		if err != nil {
			return "", fmt.Errorf("token from word: %w", err)
		}
		if token == nil {
			continue
		}
		return string(token), nil
	}

	return "", nil
}

func (t *tokenizer) TokenFromWord(w []byte) ([]byte, error) {
	for _, p := range t.patterns {
		loc := p.FindIndex(w)
		if loc == nil {
			continue
		}

		token := w[loc[0]:loc[1]]

		if before := w[:loc[0]]; len(before) != 0 {
			t.queue = append(t.queue, before)
		}
		if after := w[loc[1]:]; len(after) != 0 {
			t.queue = append(t.queue, after)
		}

		return token, nil
	}
	return nil, nil
}

// TokenPatterns returns regexps for all allowed complex token patterns.
func tokenPatterns() []regexp.Regexp {
	var out []regexp.Regexp
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		out = append(out, *re)
	}
	return out
}
