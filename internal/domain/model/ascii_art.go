package model

import "fmt"

type ASCIIArt struct {
	lines []string
}

func NewASCIIArt(lines []string) (*ASCIIArt, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("アスキーアートが空です")
	}
	return &ASCIIArt{
		lines: lines,
	}, nil
}

func (aa *ASCIIArt) Lines() []string {
	return aa.lines
}

func (aa *ASCIIArt) LineCount() int {
	return len(aa.lines)
}

func (aa *ASCIIArt) GetLine(index int) string {
	if index < 0 || index >= len(aa.lines) {
		return ""
	}
	return aa.lines[index]
}

func (aa *ASCIIArt) GetLineBySeq(seq int) string {
	if len(aa.lines) == 0 {
		return ""
	}
	return aa.lines[seq%len(aa.lines)]
}
