package token

import (
	"fmt"
	"time"

	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtClaims struct {
	jwt.RegisteredClaims
	UserID       string  `json:"uid"`
	Username     string  `json:"username"`
	Role         string  `json:"role"`
	WajibPajakID *string `json:"wpid,omitempty"`
}

type jwtProvider struct {
	secret   []byte
	ttlHours int
}

func NewJWTProvider(secret string, ttlHours int) port.TokenProvider {
	return &jwtProvider{secret: []byte(secret), ttlHours: ttlHours}
}

func (p *jwtProvider) Generate(user *entity.User) (string, int64, error) {
	expiresAt := time.Now().Add(time.Duration(p.ttlHours) * time.Hour)
	expiresIn := int64(p.ttlHours * 3600)

	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID:   user.ID.String(),
		Username: user.Username,
		Role:     string(user.Role),
	}
	if user.WajibPajakID != nil {
		s := user.WajibPajakID.String()
		claims.WajibPajakID = &s
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(p.secret)
	if err != nil {
		return "", 0, fmt.Errorf("sign token: %w", err)
	}
	return signed, expiresIn, nil
}

func (p *jwtProvider) Validate(tokenStr string) (*port.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return p.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	c, ok := token.Claims.(*jwtClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	userID, err := uuid.Parse(c.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id in token")
	}

	claims := &port.Claims{
		UserID:   userID,
		Username: c.Username,
		Role:     entity.Role(c.Role),
	}
	if c.WajibPajakID != nil {
		wpID, err := uuid.Parse(*c.WajibPajakID)
		if err == nil {
			claims.WajibPajakID = &wpID
		}
	}
	return claims, nil
}
