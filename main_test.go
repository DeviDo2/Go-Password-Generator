package main

import "testing"

func TestSelectedCharset(t *testing.T) {
	req := generateRequest{
		Lowercase: true,
		Numbers:   true,
	}

	got := selectedCharset(req)
	want := lowercaseChars + numberChars

	if got != want {
		t.Fatalf("selectedCharset() = %q, want %q", got, want)
	}
}

func TestCalculateEntropy(t *testing.T) {
	got := calculateEntropy(16, 62)
	want := 95.267

	if got < want-0.001 || got > want+0.001 {
		t.Fatalf("calculateEntropy() = %f, want near %f", got, want)
	}
}

func TestGeneratePasswordsRequiresCharset(t *testing.T) {
	_, _, _, err := generatePasswords(generateRequest{Length: 12, Count: 1})
	if err == nil {
		t.Fatal("generatePasswords() expected an error when no character set is selected")
	}
}

func TestGeneratePasswordsNormalizesLimits(t *testing.T) {
	results, _, _, err := generatePasswords(generateRequest{
		Length:    500,
		Count:     100,
		Lowercase: true,
	})
	if err != nil {
		t.Fatalf("generatePasswords() returned error: %v", err)
	}

	if len(results) != maxCount {
		t.Fatalf("len(results) = %d, want %d", len(results), maxCount)
	}

	if len(results[0].Value) != maxLength {
		t.Fatalf("password length = %d, want %d", len(results[0].Value), maxLength)
	}
}
