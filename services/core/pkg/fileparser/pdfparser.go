package fileparser

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type PDFParser struct {
	PdftotextBin string
	PdfimagesBin string
}

func (p *PDFParser) Parse(r io.Reader) (*ParsedContent, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading pdf: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "pdfparse")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	pdfPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(pdfPath, data, 0644); err != nil {
		return nil, err
	}

	textCmd := exec.Command(p.PdftotextBin, "-layout", pdfPath, "-")
	textOut, err := textCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("pdftotext: %w", err)
	}
	fullText := string(textOut)

	pagesText := strings.Split(fullText, "\f")

	imgDir := filepath.Join(tmpDir, "images")
	if err := os.Mkdir(imgDir, 0755); err != nil {
		return nil, err
	}
	imgPrefix := filepath.Join(imgDir, "img")
	imgCmd := exec.Command(p.PdfimagesBin, "-j", "-p", pdfPath, imgPrefix)
	if out, err := imgCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("pdfimages: %w output: %s", err, string(out))
	}

	imgFiles, _ := filepath.Glob(imgPrefix + "*")
	sort.Strings(imgFiles)

	var imagesByPage = make(map[int][]string)
	globalImages := make([]ImagePart, 0)
	for _, fpath := range imgFiles {
		base := filepath.Base(fpath)
		parts := strings.SplitN(base, "-", 3)
		if len(parts) < 2 {
			continue
		}

		pageNum, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		imgData, err := os.ReadFile(fpath)
		if err != nil {
			continue
		}

		compressedData, err := CompressImage(imgData, maxImageSize, jpegQuality)
		if err != nil {
			continue
		}

		imagesByPage[pageNum] = append(imagesByPage[pageNum], base)
		globalImages = append(globalImages, ImagePart{ID: base, Data: compressedData})
	}

	var resultText strings.Builder

	for i, pageText := range pagesText {
		pageNum := i + 1
		pagePrifix := fmt.Sprintf("Страница %d:\n", pageNum)
		resultText.WriteString(pagePrifix)
		resultText.WriteString(pageText)
		if imgs, ok := imagesByPage[pageNum]; ok {
			resultText.WriteString("Страница содержит изображения:")
			for _, imgID := range imgs {
				text := fmt.Sprintf(" %s", imgID)
				resultText.WriteString(text)
			}
			resultText.WriteString("\n")
		}
	}

	return &ParsedContent{
		Text:   resultText.String(),
		Images: globalImages,
	}, nil
}
