package jwt

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
)

// New func to create a new instance of JWT
func New(signingKey []byte, ttl int64) *JWT {
	return &JWT{signingKey, ttl}
}

type (
	// JWT structure
	JWT struct {
		signingKey []byte
		ttl        int64
	}

	// Claims is a custom JWT claims structure
	Claims struct {
		jwt.StandardClaims
		User map[string]interface{} `json:"user"`
	}

	// Response structure
	Response struct {
		AccessToken string `json:"access_token"`
		ExpiresAt   int64  `json:"expires_at"`
	}
)

// SigningKey returns JWT signing key
func (j *JWT) SigningKey() []byte {
	return j.signingKey
}

// TTL returns JWT life time in seconds
func (j *JWT) TTL() int64 {
	return j.ttl
}

// NewWithPayload function creates a new token response with payload
func (j *JWT) NewWithPayload(payload map[string]interface{}) (*Response, error) {
	return j.new(j.NewClaims(payload))
}

// NewWithClaims function creates a new token response with claims
// Alias of the private function JWT.new()
func (j *JWT) NewWithClaims(claims *Claims) (*Response, error) {
	return j.new(claims)
}

// NewClaims returns new JWTClaims structure instance with populated default fields
func (j *JWT) NewClaims(payload map[string]interface{}) *Claims {
	claims := &Claims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(j.ttl) * time.Second).Unix(),
			Id:        uuid.NewV1().String(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		}, payload,
	}
	return claims
}

// helper to generate new JWT
func (j *JWT) new(claims *Claims) (*Response, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.signingKey)
	if err != nil {
		return nil, err
	}
	return &Response{
		AccessToken: token,
		ExpiresAt:   claims.ExpiresAt,
	}, nil
}
