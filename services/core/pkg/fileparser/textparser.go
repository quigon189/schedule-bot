package fileparser

import (
	"fmt"
	"io"
)

type TextParser struct{}

func (p *TextParser) Parse(r io.Reader) (*ParsedContent, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	return &ParsedContent{
		Text: string(content),
	}, nil
}
