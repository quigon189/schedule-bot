package fileparser

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gabriel-vasile/mimetype"
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

func Parse(r io.Reader) (*ParsedContent, error) {
	header := bytes.NewBuffer(nil)
	mtype, err := mimetype.DetectReader(io.TeeReader(r, header))
	if err != nil {
		return nil, fmt.Errorf("detect file type: %w", err)
	}
	p, ok := parsers[mtype.String()]
	if !ok {
		return nil, fmt.Errorf("unsupported format: %s", mtype.String())
	}

	return p.Parse(io.MultiReader(header, r))
}

func init() {
	pdfParser := &PDFParser{
		PdftotextBin: "pdftotext",
		PdfimagesBin: "pdfimages",
	}
	Register("application/pdf", pdfParser)

	Register("text/plain; charset=utf-8", &TextParser{})

	imageParser := &ImageParser{}
	Register("image/jpeg", imageParser)
	Register("image/png", imageParser)
	Register("image/gif", imageParser)
	Register("image/bmp", imageParser)
}
