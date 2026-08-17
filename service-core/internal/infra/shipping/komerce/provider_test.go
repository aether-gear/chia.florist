package komerce

import (
	"testing"
)

func TestNormalizeCourierCode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"SPX Standard", "spx"},
		{"Shopee Express", "spx"},
		{"spx", "spx"},
		{"JNE Express REG", "jne"},
		{"J&T Express", "jnt"},
		{"JNT EZ", "jnt"},
		{"SiCepat REG", "sicepat"},
		{"POS Indonesia", "pos"},
		{"TIKI ONS", "tiki"},
		{"Ninja Express", "ninja"},
		{"IDExpress", "ide"},
		{"Lion Parcel", "lion"},
		{"Wahana Express", "wahana"},
		{"custom_code", "custom_code"},
	}

	for _, tt := range tests {
		got := normalizeCourierCode(tt.input)
		if got != tt.expected {
			t.Errorf("normalizeCourierCode(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}
