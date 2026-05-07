package main

import (
	"crypto/rand"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"math"
	"math/big"
	"net/http"
	"strings"
	"time"
)

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars    = "0123456789"
	symbolChars    = "!@#$%^&*()-_=+[]{};:,.<>?/|~"

	defaultLength = 16
	defaultCount  = 4
	maxLength     = 128
	maxCount      = 50
)

//go:embed static/*
var embeddedFiles embed.FS

type generateRequest struct {
	Length    int  `json:"length"`
	Count     int  `json:"count"`
	Lowercase bool `json:"lowercase"`
	Uppercase bool `json:"uppercase"`
	Numbers   bool `json:"numbers"`
	Symbols   bool `json:"symbols"`
}

type passwordResult struct {
	Value       string  `json:"value"`
	EntropyBits float64 `json:"entropyBits"`
}

type generateResponse struct {
	Passwords     []passwordResult `json:"passwords"`
	CharacterPool int              `json:"characterPool"`
	EntropyBits   float64          `json:"entropyBits"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func main() {
	staticFS, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/generate", generateHandler)
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           securityHeaders(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Password generator running at http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Request body must be valid JSON."})
		return
	}

	results, poolSize, entropy, err := generatePasswords(req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, generateResponse{
		Passwords:     results,
		CharacterPool: poolSize,
		EntropyBits:   entropy,
	})
}

func generatePasswords(req generateRequest) ([]passwordResult, int, float64, error) {
	req = normalizeRequest(req)
	charset := selectedCharset(req)
	if charset == "" {
		return nil, 0, 0, errors.New("Select at least one character set.")
	}

	entropy := calculateEntropy(req.Length, len(charset))
	results := make([]passwordResult, req.Count)
	for i := range results {
		password, err := randomPassword(req.Length, charset)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("could not generate password: %w", err)
		}

		results[i] = passwordResult{
			Value:       password,
			EntropyBits: roundEntropy(entropy),
		}
	}

	return results, len(charset), roundEntropy(entropy), nil
}

func normalizeRequest(req generateRequest) generateRequest {
	if req.Length <= 0 {
		req.Length = defaultLength
	}
	if req.Length > maxLength {
		req.Length = maxLength
	}

	if req.Count <= 0 {
		req.Count = defaultCount
	}
	if req.Count > maxCount {
		req.Count = maxCount
	}

	return req
}

func selectedCharset(req generateRequest) string {
	var builder strings.Builder
	if req.Lowercase {
		builder.WriteString(lowercaseChars)
	}
	if req.Uppercase {
		builder.WriteString(uppercaseChars)
	}
	if req.Numbers {
		builder.WriteString(numberChars)
	}
	if req.Symbols {
		builder.WriteString(symbolChars)
	}
	return builder.String()
}

func randomPassword(length int, charset string) (string, error) {
	var builder strings.Builder
	builder.Grow(length)

	maxIndex := big.NewInt(int64(len(charset)))
	for builder.Len() < length {
		index, err := rand.Int(rand.Reader, maxIndex)
		if err != nil {
			return "", err
		}
		builder.WriteByte(charset[index.Int64()])
	}

	return builder.String(), nil
}

func calculateEntropy(length int, poolSize int) float64 {
	if length <= 0 || poolSize <= 1 {
		return 0
	}
	return float64(length) * math.Log2(float64(poolSize))
}

func roundEntropy(value float64) float64 {
	return math.Round(value*10) / 10
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("could not write JSON response: %v", err)
	}
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; base-uri 'none'; form-action 'none'")
		next.ServeHTTP(w, r)
	})
}
