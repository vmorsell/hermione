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
	NextToken() (string, error)
	HasMoreTokens() bool
}

type tokenizer struct {
	r          *bufio.Reader
	tokenQueue []string
	patterns   []regexp.Regexp
}

func NewTokenizer(reader io.Reader) Tokenizer {
	return &tokenizer{
		r:        bufio.NewReader(reader),
		patterns: tokenPatterns(),
	}
}

// NextToken reads the next token from the reader and returns it as a string.
func (t *tokenizer) NextToken() (string, error) {
	for len(t.tokenQueue) == 0 {
		word, err := t.ReadNextWord()
		if err != nil {
			return "", fmt.Errorf("read next word: %w", err)
		}

		tokens, err := t.TokensFromWord(word)
		if err != nil {
			return "", fmt.Errorf("tokens from word: %w", err)
		}

		t.tokenQueue = append(t.tokenQueue, tokens...)
	}

	next := t.tokenQueue[0]
	if len(t.tokenQueue) > 1 {
		t.tokenQueue = t.tokenQueue[1:]
	} else {
		t.tokenQueue = nil
	}
	return next, nil
}

func (t *tokenizer) ReadNextWord() (string, error) {
	var out bytes.Buffer

readByte:
	for {
		b, err := t.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				// End word at EOF.
				break
			}
			return "", fmt.Errorf("read byte: %w", err)
		}

		for _, s := range stopBytes {
			if s == b {
				break readByte
			}
		}

		c := bytes.ToLower([]byte{b})[0]
		out.WriteByte(c)
	}

	return out.String(), nil
}

func (t *tokenizer) TokensFromWord(w string) ([]string, error) {
	var out []string
	for _, p := range t.patterns {
		if tokens := p.FindAllString(w, -1); tokens != nil {
			out = append(out, tokens...)
			w = p.ReplaceAllString(w, "")
		}
	}
	return out, nil
}

// HasMoreTokens returns if there are more tokens to read.
func (t *tokenizer) HasMoreTokens() bool {
	_, err := t.r.Peek(1)
	return err == nil
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
