package postgres

import (
	"context"
	"errors"

	"github.com/bpka/tps-pkb/internal/domain/dto"
	"github.com/bpka/tps-pkb/internal/domain/entity"
	"github.com/bpka/tps-pkb/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct{ pool *pgxpool.Pool }

func NewUserRepository(pool *pgxpool.Pool) port.UserRepository {
	return &userRepo{pool: pool}
}

func (r *userRepo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, role, wajib_pajak_id, created_at, updated_at
		 FROM users WHERE username = $1`, username)
	return scanUser(row)
}

func (r *userRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	row := r.pool.QueryRow(ctx,
		`SELECT id, username, password_hash, role, wajib_pajak_id, created_at, updated_at
		 FROM users WHERE id = $1`, id)
	return scanUser(row)
}

func (r *userRepo) Save(ctx context.Context, user *entity.User) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, username, password_hash, role, wajib_pajak_id)
		 VALUES ($1, $2, $3, $4, $5)`,
		user.ID, user.Username, user.PasswordHash, user.Role, user.WajibPajakID)
	return err
}

func (r *userRepo) Update(ctx context.Context, user *entity.User) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET password_hash=$1, role=$2, wajib_pajak_id=$3, updated_at=NOW()
		 WHERE id=$4`,
		user.PasswordHash, user.Role, user.WajibPajakID, user.ID)
	return err
}

func scanUser(row pgx.Row) (*entity.User, error) {
	var u entity.User
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.WajibPajakID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dto.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
