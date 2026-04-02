package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type TokenService struct {
	secret string
}

func NewTokenService(secret string) *TokenService {
	return &TokenService{secret: secret}
}

// Generate crée un token signé encodant le compteID.
func (t *TokenService) Generate(compteID uint) string {
	payload := base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%d", compteID)))
	mac := hmac.New(sha256.New, []byte(t.secret))
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))
	return payload + "." + sig
}

// Validate vérifie la signature et retourne le compteID.
func (t *TokenService) Validate(token string) (uint, error) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return 0, errors.New("token invalide")
	}

	mac := hmac.New(sha256.New, []byte(t.secret))
	mac.Write([]byte(parts[0]))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(parts[1]), []byte(expected)) {
		return 0, errors.New("token invalide")
	}

	payloadBytes, err := base64.URLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, errors.New("token invalide")
	}

	id, err := strconv.ParseUint(string(payloadBytes), 10, 64)
	if err != nil {
		return 0, errors.New("token invalide")
	}

	return uint(id), nil
}
