package session

import (
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	Address common.Address `json:"address"`
}

type Manager struct {
	secretKey []byte
	expiry    time.Duration
}

var (
	sessionManager *Manager
	initOnce       sync.Once
)

func Init(secretKey []byte, expiry time.Duration) error {
	if sessionManager != nil {
		return fmt.Errorf("session manager already initialized")
	}

	initOnce.Do(func() {
		sessionManager = &Manager{
			secretKey: secretKey,
			expiry:    expiry,
		}
	})

	return nil
}

func CreateToken(address common.Address) (string, error) {
	claims := Claims{
		Address: address,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(sessionManager.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(sessionManager.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*common.Address, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return sessionManager.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return &claims.Address, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
