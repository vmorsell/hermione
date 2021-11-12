package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

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

out:
	for {
		b, err := t.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				// Token can end on EOF.
				break
			}
			return "", fmt.Errorf("read byte: %w", err)
		}

		// Token can also end on a space or a line break character.
		if b == ' ' || b == '\n' {
			break
		}

		for _, f := range filter {
			if b == f {
				continue out
			}
		}

		lower := bytes.ToLower([]byte{b})
		token.WriteByte(lower[0])
	}
	return token.String(), nil
}

// HasMoreTokens returns if there are more tokens to read.
func (t *tokenizer) HasMoreTokens() bool {
	_, err := t.r.Peek(1)
	return err == nil
}
