package komerce

import (
	"testing"
)

func TestNormalizeCourierCode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"JNE Express REG", "jne"},
		{"J&T Express", "jnt"},
		{"JNT EZ", "jnt"},
		{"Ninja Express", "ninja"},
		{"TIKI ONS", "tiki"},
		{"POS Indonesia", "pos"},
		{"AnterAja Regular", "anteraja"},
		{"SAP Express", "sap"},
		{"Lion Parcel", "lion"},
		{"Wahana Express", "wahana"},
		{"First Logistics", "first"},
		{"IDExpress", "ide"},
	}

	for _, tt := range tests {
		got := normalizeCourierCode(tt.input)
		if got != tt.expected {
			t.Errorf("normalizeCourierCode(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

func TestIsValidKomerceCourier(t *testing.T) {
	validCouriers := []string{
		"jne", "JNE Express",
		"jnt", "J&T Express",
		"ninja", "Ninja Express",
		"tiki", "TIKI ONS",
		"pos", "POS Indonesia",
		"anteraja", "AnterAja Eco",
		"sap", "SAP Express",
		"lion", "Lion Parcel",
		"wahana", "Wahana",
		"first", "First Logistics",
		"ide", "IDExpress",
	}

	for _, c := range validCouriers {
		if !isValidKomerceCourier(c) {
			t.Errorf("isValidKomerceCourier(%q) = false; want true", c)
		}
	}

	invalidCouriers := []string{
		"spx", "Shopee Express", "sicepat", "sentral", "rex", "custom_code", "",
	}

	for _, c := range invalidCouriers {
		if isValidKomerceCourier(c) {
			t.Errorf("isValidKomerceCourier(%q) = true; want false", c)
		}
	}
}
