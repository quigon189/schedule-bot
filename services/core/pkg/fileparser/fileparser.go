package fileparser

import (
	"fmt"
	"io"
	"path/filepath"
)

type ImagePart struct {
	ID string
	Data []byte
}

type ParsedContent struct {
	Text string
	Images []ImagePart
}

type Parser interface {
	Parse(r io.Reader) (*ParsedContent, error)
}

var parsers = map[string]Parser{}

func Register(ext string, p Parser) {
	parsers[ext] = p
}

func Parse(r io.Reader, filename string) (*ParsedContent, error) {
	ext := filepath.Ext(filename)
	p, ok := parsers[ext]
	if !ok {
		return nil, fmt.Errorf("unsupported format: %s", ext)
	}

	return p.Parse(r)
}

func init() {
	pdfParser := &PDFParser{
		PdftotextBin: "pdftotext",
		PdfimagesBin: "pdfimages",
	}
	Register(".pdf", pdfParser)
}
