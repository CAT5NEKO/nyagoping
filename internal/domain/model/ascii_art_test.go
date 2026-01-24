package model

import (
	"testing"
)

func TestNewASCIIArt(t *testing.T) {
	tests := []struct {
		name    string
		lines   []string
		wantErr bool
	}{
		{
			name:    "有効なアスキーアート",
			lines:   []string{"line1", "line2", "line3"},
			wantErr: false,
		},
		{
			name:    "空のアスキーアート",
			lines:   []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			art, err := NewASCIIArt(tt.lines)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewASCIIArt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && art.LineCount() != len(tt.lines) {
				t.Errorf("NewASCIIArt() LineCount = %v, want %v", art.LineCount(), len(tt.lines))
			}
		})
	}
}

func TestASCIIArt_GetLineBySeq(t *testing.T) {
	lines := []string{"line1", "line2", "line3"}
	art, _ := NewASCIIArt(lines)

	tests := []struct {
		name string
		seq  int
		want string
	}{
		{
			name: "最初の行",
			seq:  0,
			want: "line1",
		},
		{
			name: "2番目の行",
			seq:  1,
			want: "line2",
		},
		{
			name: "循環（行数を超える）",
			seq:  3,
			want: "line1",
		},
		{
			name: "循環（大きい数）",
			seq:  7,
			want: "line2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := art.GetLineBySeq(tt.seq); got != tt.want {
				t.Errorf("GetLineBySeq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestASCIIArt_LineCount(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  int
	}{
		{
			name:  "3行",
			lines: []string{"a", "b", "c"},
			want:  3,
		},
		{
			name:  "1行",
			lines: []string{"single"},
			want:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			art, _ := NewASCIIArt(tt.lines)
			if got := art.LineCount(); got != tt.want {
				t.Errorf("LineCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
