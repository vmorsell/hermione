package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
)

var patterns = []string{
	`https?:\/\/\S+`,                     // Web addresses
	`[\w\._-]+@[\w\._-]+\.\w+`,           // Email addresses
	`\d+/\d+`,                            // Twenty-four seven
	`\w+'\w+`,                            // Contraditions
	`'em`,                                // 'em
	`\$\d+`,                              // UDS amounts
	`a\.m\.|p\.m\.|a\.d\.|b\.c\.|e\.g\.`, // Abbreviations
	`(\w+[\.,-])+\w+`,                    // Chars or digits separated with dot, comma or dash
	`(\d{2}:)+\d{2}`,                     // Time
	`\+\d+`,                              // Numbers beginning with +
}

type Tokenizer interface {
	NextToken() (string, error)
	HasMoreTokens() bool
}

type tokenizer struct {
	r *bufio.Reader
}

func NewTokenizer(reader io.Reader) Tokenizer {
	return &tokenizer{
		r: bufio.NewReader(reader),
	}
}

// NextToken reads the next token from the reader and returns it as a string.
func (t *tokenizer) NextToken() (string, error) {
	var token bytes.Buffer

	filter := []byte("!\"#€%&/()=?´`'©@£$∞§|[]≈±¨^*,.-;:_<>§°¶")

	for {
		b, err := t.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				// Token can end on EOF.
				break
			}
			return "", fmt.Errorf("read byte: %w", err)
		}

		// Spaces can end a token, but only if we already have bytes in the buffer.
		if b == ' ' {
			if token.Len() == 0 {
				continue
			}
			break
		}

		// Always end tokens on line breaks.
		if b == '\n' {
			break
		}

		lower := bytes.ToLower([]byte{b})
		token.WriteByte(lower[0])
	}

	patterns, err := TokenPatterns()
	if err != nil {
		return "", fmt.Errorf("token patterns: %w", err)
	}
	for _, p := range patterns {
		if p.Match(token.Bytes()) {
			return token.String(), nil
		}
	}

	// Token does not match pattern. Remove special characters.
	var out bytes.Buffer
nextByte:
	for _, b := range token.Bytes() {
		for _, f := range filter {
			if b == f {
				continue nextByte
			}
		}
		out.WriteByte(b)
	}

	// Ignore empty tokens.
	if out.Len() == 0 {
		return t.NextToken()
	}

	return out.String(), nil
}

// HasMoreTokens returns if there are more tokens to read.
func (t *tokenizer) HasMoreTokens() bool {
	_, err := t.r.Peek(1)
	return err == nil
}

// TokenPatterns returns regexps for all allowed complex token patterns.
func TokenPatterns() ([]regexp.Regexp, error) {
	var out []regexp.Regexp
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("compile pattern `%s`: %w", p, err)
		}
		out = append(out, *re)
	}
	return out, nil
}
