package fileparser

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
)

const maxImageSize = 512
const jpegQuality = 80

type ImageParser struct {}

func (p *ImageParser) Parse(r io.Reader) (*ParsedContent, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	_, _, err = image.DecodeConfig(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	comressedData, err := CompressImage(data, maxImageSize, jpegQuality)
	if err != nil {
		return nil, err
	}

	imagePart := ImagePart{
		ID: "image.jpeg",
		Data: comressedData,
	}

	return &ParsedContent{
		Text: "",
		Images: []ImagePart{imagePart},
	}, nil
}

func CompressImage(data []byte, maxSize int, quality int) ([]byte, error) {
	img, _, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	if width > height && width > maxSize {
		height = height * maxSize / width
		width = maxSize
	} else if height > maxSize {
		width = width * maxSize / height
		height = maxSize
	}

	resized := imaging.Resize(img, width, height, imaging.Lanczos)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resized, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
