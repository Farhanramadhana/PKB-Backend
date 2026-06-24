package port

import (
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/google/uuid"
)

type Claims struct {
	UserID       uuid.UUID
	Username     string
	Role         entity.Role
	WajibPajakID *uuid.UUID
}

//go:generate mockgen -source=token_provider.go -destination=../../mock/mock_token_provider.go -package=mock
type TokenProvider interface {
	Generate(user *entity.User) (token string, expiresIn int64, err error)
	Validate(token string) (*Claims, error)
}
