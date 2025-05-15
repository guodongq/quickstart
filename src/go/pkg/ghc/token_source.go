package ghc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"golang.org/x/oauth2"
	"sync"
	"time"
)

const (
	JWTHeaderType         = "JWT"
	JWTSigningAlgorithm   = "RS256"
	TokenExpirationMargin = 30 * time.Second
	DefaultTokenValidity  = 10 * time.Minute
	TokenBackdateDuration = 60 * time.Second
)

type JWTGenerator struct {
	clientID      string
	privateKeyPEM []byte
}

func NewJWTGenerator(clientID string, privateKeyPEM []byte) *JWTGenerator {
	return &JWTGenerator{
		clientID:      clientID,
		privateKeyPEM: privateKeyPEM,
	}
}

func (g *JWTGenerator) Generate() (string, time.Time, error) {
	now := time.Now().UTC()
	claims := jwtClaims{
		Issuer:    g.clientID,
		IssuedAt:  now.Add(-TokenBackdateDuration).Unix(),
		ExpiresAt: now.Add(DefaultTokenValidity).Unix(),
	}

	headerB64, err := encodeJWTHeader()
	if err != nil {
		return "", time.Time{}, err
	}

	payloadB64, err := encodeClaims(claims)
	if err != nil {
		return "", time.Time{}, err
	}

	signature, err := g.sign(headerB64, payloadB64)
	if err != nil {
		return "", time.Time{}, err
	}
	return fmt.Sprintf("%s.%s.%s", headerB64, payloadB64, signature), now.Add(DefaultTokenValidity), nil
}

func (g *JWTGenerator) sign(headerB64, payloadB64 string) (string, error) {
	block, _ := pem.Decode(g.privateKeyPEM)
	if block == nil {
		return "", fmt.Errorf("invalid PEM format: no PEM blocks found")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("parse private key failed: %w", err)
	}

	signingInput := headerB64 + "." + payloadB64
	hashed := sha256.Sum256([]byte(signingInput))

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", fmt.Errorf("signing failed: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(signature), nil
}

func encodeClaims(claims jwtClaims) (string, error) {
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal payload failed: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(payloadJSON), nil
}

func encodeJWTHeader() (string, error) {
	header := struct {
		Type string `json:"typ"`
		Alg  string `json:"alg"`
	}{
		Type: JWTHeaderType,
		Alg:  JWTSigningAlgorithm,
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("marshal header failed: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(headerJSON), nil
}

type jwtClaims struct {
	Issuer    string `json:"iss"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

type CachedTokenSource struct {
	generator *JWTGenerator
	mu        sync.RWMutex
	token     *oauth2.Token
}

func NewCachedTokenSource(generator *JWTGenerator) oauth2.TokenSource {
	return &CachedTokenSource{
		generator: generator,
	}
}

func (c *CachedTokenSource) Token() (*oauth2.Token, error) {
	c.mu.RLock()
	if c.tokenValid() {
		defer c.mu.RUnlock()
		return c.token, nil
	}
	c.mu.RUnlock()

	tokenString, expiry, err := c.generator.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate JWT failed: %w", err)
	}

	c.token = &oauth2.Token{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		Expiry:      expiry.Add(-TokenExpirationMargin),
	}

	return c.token, nil
}

func (c *CachedTokenSource) tokenValid() bool {
	return c.token != nil && time.Now().Before(c.token.Expiry)
}
