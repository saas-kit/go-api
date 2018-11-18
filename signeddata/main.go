package signeddata

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	expirationTimeKey string = "exp"
)

// Predefined errors
var (
	ErrCouldNotDecodeString = errors.New("could not decode signed string")
	ErrExpired              = errors.New("signed string was expired")
)

// Encode data
func Encode(key string, payload map[string]interface{}, exp time.Time) (string, error) {
	payload[expirationTimeKey] = exp.UTC().Unix()
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	encodedPayload := encodeSegment(jsonPayload)
	hash, err := bcrypt.GenerateFromPassword([]byte(string(encodedPayload)+key), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", encodedPayload, encodeSegment(hash)), nil
}

// Decode data
func Decode(key, signedString string) (map[string]interface{}, error) {
	parts := strings.Split(signedString, ".")
	if len(parts) != 2 {
		return nil, ErrCouldNotDecodeString
	}
	jsonData, err := decodeSegment(parts[0])
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}
	if err := checkExpiration(data); err != nil {
		return nil, err
	}
	if err := verify(key, signedString); err != nil {
		return nil, err
	}
	return data, nil
}

// Verify signature
func verify(signingKey, signedString string) error {
	parts := strings.Split(signedString, ".")
	if len(parts) != 2 {
		return ErrCouldNotDecodeString
	}
	signature, err := decodeSegment(parts[1])
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword(signature, []byte(parts[0]+signingKey))
}

// checkExpiration helper
func checkExpiration(data map[string]interface{}) error {
	exp, ok := data[expirationTimeKey]
	if !ok {
		return ErrExpired
	}
	if expTime, ok := exp.(float64); !ok || int64(expTime) < time.Now().Unix() {
		return ErrExpired
	}
	return nil
}

// Encode segment specific base64url encoding with padding stripped
func encodeSegment(seg []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(seg), "=")
}

// Decode segment specific base64url encoding with padding stripped
func decodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}
	return base64.URLEncoding.DecodeString(seg)
}
