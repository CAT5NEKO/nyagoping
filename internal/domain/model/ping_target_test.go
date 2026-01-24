package model

import (
	"testing"
)

func TestNewPingTarget(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		wantErr bool
	}{
		{
			name:    "有効なホスト名",
			host:    "example.tld",
			wantErr: false,
		},
		{
			name:    "空のホスト名",
			host:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target, err := NewPingTarget(tt.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPingTarget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && target.Host() != tt.host {
				t.Errorf("NewPingTarget() host = %v, want %v", target.Host(), tt.host)
			}
		})
	}
}

func TestPingTarget_Host(t *testing.T) {
	host := "example.tld"
	target, _ := NewPingTarget(host)

	if got := target.Host(); got != host {
		t.Errorf("Host() = %v, want %v", got, host)
	}
}
