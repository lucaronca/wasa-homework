package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lucaronca/wasa-homework/service/database"
)

type BansRepository interface {
	// Getters
	GetBanExists(int, int) (bool, error)
	// Setters
	SetBan(int, int) error
	RemoveBan(int, int) error
	// Relations builders
	WithoutBanned(int) Relation
	WithoutBanners(int) Relation
}

type bansRepository struct {
	database.AppDatabase
}

func NewBansRepository(db database.AppDatabase) (BansRepository, error) {
	if db == nil {
		return nil, errors.New("database is required")
	}

	return &bansRepository{
		db,
	}, nil
}

func (r *bansRepository) SetBan(userId int, bannedId int) error {
	if _, err := r.Conn().Exec(`
		INSERT OR IGNORE INTO user_bans (user_id, banned_id)
		VALUES (?, ?);
	`, userId, bannedId); err != nil {
		return err
	}
	return nil
}

func (r *bansRepository) RemoveBan(userId int, bannedId int) error {
	if _, err := r.Conn().Exec(`
		DELETE FROM user_bans
		WHERE user_id=? AND banned_id=?;
	`, userId, bannedId); err != nil {
		return err
	}
	return nil
}

func (r *bansRepository) GetBanExists(userId int, targetId int) (bool, error) {
	var exists int
	err := r.Conn().QueryRow(`
		SELECT COUNT(1) FROM user_bans
		WHERE user_id=? AND banned_id=?;
	`, userId, targetId).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if exists == 0 {
		return false, nil
	}
	return true, nil
}

// Relations builders
func (r *bansRepository) WithoutBanned(userId int) Relation {
	return Relation(func(string) string {
		return fmt.Sprintf(
			"WHERE user_id NOT IN (SELECT banned_id from user_bans WHERE user_id = %d)",
			userId,
		)
	})
}

func (r *bansRepository) WithoutBanners(userId int) Relation {
	return Relation(func(entity string) string {
		var userIdField string
		if entity == "user" {
			userIdField = "id"
		} else {
			userIdField = "user_id"
		}
		return fmt.Sprintf(
			"WHERE %s NOT IN (SELECT user_id from user_bans WHERE banned_id = %d)",
			userIdField,
			userId,
		)
	})
}
