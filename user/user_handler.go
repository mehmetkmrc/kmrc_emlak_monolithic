package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func (r *UserRepository) UpdateUser(
	ctx context.Context,
	userID string,
	name string,
	surname string,
	email string,
	phone string,
) error {

	query := `
        UPDATE users
        SET 
            first_name = $1,
            last_name = $2,
            email = $3,
            phone = $4,
            updated_at = NOW()
        WHERE user_id = $5
    `

	_, err := r.db.Exec(ctx, query, name, surname, email, phone, userID)
	if err != nil {
		return err
	}

	return nil
}
