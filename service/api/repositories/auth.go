package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lucaronca/wasa-homework/service/database"
)

type AuthRepository interface {
	// Getters
	GetToken(relations ...Relation) (string, error)
	// Setters
	SetToken(userId int, token string) error
	// Relations builders
	WithTokens() Relation
	FilterByToken(token string) Relation
}

type authRepository struct {
	database.AppDatabase
}

func NewAuthRepository(db database.AppDatabase) (AuthRepository, error) {
	if db == nil {
		return nil, errors.New("database is required")
	}

	return &authRepository{
		db,
	}, nil
}

func (r *authRepository) GetToken(relations ...Relation) (string, error) {
	q := queryBuilder("user_token", relations...)
	var token string
	err := r.Conn().QueryRow(
		fmt.Sprintf(`
		SELECT token FROM user_tokens
		%s
	`, q)).Scan(&token)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *authRepository) SetToken(userId int, token string) error {
	if _, err := r.Conn().Exec(`
		INSERT INTO user_tokens (token, user_id)
		VALUES (?, ?)
	`, token, userId); err != nil {
		return err
	}
	return nil
}

// Relations builders
func (r *authRepository) WithTokens() Relation {
	return Relation(func(entity string) string {
		return fmt.Sprintf(
			"INNER JOIN user_tokens ON %ss.id = %s_id",
			entity,
			entity,
		)
	})
}

func (r *authRepository) FilterByToken(token string) Relation {
	return Relation(func(entity string) string {
		return fmt.Sprintf(
			"WHERE token=\"%s\"",
			token,
		)
	})
}
