package model

import (
	"testing"
)

func TestNewPingConfig(t *testing.T) {
	tests := []struct {
		name       string
		count      int
		privileged bool
		wantErr    bool
	}{
		{
			name:       "有効な設定",
			count:      10,
			privileged: false,
			wantErr:    false,
		},
		{
			name:       "0回のPing",
			count:      0,
			privileged: true,
			wantErr:    false,
		},
		{
			name:       "負のカウント",
			count:      -1,
			privileged: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewPingConfig(tt.count, tt.privileged)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPingConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if config.Count() != tt.count {
					t.Errorf("Count() = %v, want %v", config.Count(), tt.count)
				}
				if config.Privileged() != tt.privileged {
					t.Errorf("Privileged() = %v, want %v", config.Privileged(), tt.privileged)
				}
			}
		})
	}
}

func TestPingConfig_SetCount(t *testing.T) {
	config, _ := NewPingConfig(10, false)

	tests := []struct {
		name     string
		newCount int
		wantErr  bool
	}{
		{
			name:     "有効な変更",
			newCount: 20,
			wantErr:  false,
		},
		{
			name:     "0への変更",
			newCount: 0,
			wantErr:  false,
		},
		{
			name:     "負の値への変更",
			newCount: -5,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.SetCount(tt.newCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && config.Count() != tt.newCount {
				t.Errorf("Count() after SetCount = %v, want %v", config.Count(), tt.newCount)
			}
		})
	}
}
