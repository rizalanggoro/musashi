package utils

import (
	"testing"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    uint
		wantErr bool
	}{
		{
			name:    "valid positive integer",
			input:   "123",
			want:    123,
			wantErr: false,
		},
		{
			name:    "valid small integer",
			input:   "1",
			want:    1,
			wantErr: false,
		},
		{
			name:    "zero is parsed as 0",
			input:   "0",
			want:    0,
			wantErr: false,
		},
		{
			name:    "negative returns error",
			input:   "-5",
			want:    0,
			wantErr: true,
		},
		{
			name:    "non-numeric string",
			input:   "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseID(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseID(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("ParseID(%q) = %d, want %d", tc.input, got, tc.want)
			}
		})
	}
}

func TestFormatID(t *testing.T) {
	if got := FormatID(42); got != "42" {
		t.Errorf("FormatID(42) = %s, want '42'", got)
	}
	if got := FormatID(0); got != "0" {
		t.Errorf("FormatID(0) = %s, want '0'", got)
	}
}
