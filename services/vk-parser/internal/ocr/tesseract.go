package ocr

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Tesseract struct {
	path string
	lang string
}

func NewTesseract(path string) *Tesseract {
	return &Tesseract{
		path: path,
		lang: "rus+eng",
	}
}

func (t *Tesseract) ExtractTextFromImage(imageBytes []byte) (string, error) {
	tmpFile, err := os.CreateTemp("", "ocr-*.jpg")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(imageBytes); err != nil {
		return "", err
	}

	cmd := exec.Command(t.path, tmpFile.Name(), "stdout", "-l", t.lang)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("tesseract error: %v, output: %s", err, out.String())
	}

	return out.String(), nil
}
