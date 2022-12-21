package repositories

import (
	"errors"
	"fmt"

	"github.com/lucaronca/wasa-homework/service/database"
)

type FollowsRepository interface {
	// Setters
	SetFollow(int, int) error
	RemoveFollow(int, int) error
	// Relation builders
	FilterByFollowerId(int) Relation
	FilterByFollowingId(int) Relation
	WithTotalFollowings() Relation
	WithTotalFollowers() Relation
}

type followsRepository struct {
	database.AppDatabase
}

func NewFollowsRepository(db database.AppDatabase) (FollowsRepository, error) {
	if db == nil {
		return nil, errors.New("database is required")
	}

	return &followsRepository{
		db,
	}, nil
}

func (r *followsRepository) SetFollow(followerId int, followingId int) error {
	if _, err := r.Conn().Exec(`
		INSERT OR IGNORE INTO follows (follower_id, following_id)
		VALUES (?, ?)
	`, followerId, followingId); err != nil {
		return err
	}
	return nil
}

func (r *followsRepository) RemoveFollow(followerId int, followingId int) error {
	if _, err := r.Conn().Exec(`
		DELETE FROM follows
		WHERE follower_id=? AND following_id=?
	`, followerId, followingId); err != nil {
		return err
	}
	return nil
}

// Relations builders
func (r *followsRepository) FilterByFollowerId(followerId int) Relation {
	return Relation(func(entity string) string {
		var userIdField string
		if entity == "user" {
			userIdField = "id"
		} else {
			userIdField = "user_id"
		}
		return fmt.Sprintf(`
				INNER JOIN (
					SELECT following_id FROM follows
					WHERE follower_id=%d
				) ON following_id = %ss.%s
			`,
			followerId,
			entity,
			userIdField,
		)
	})
}

func (r *followsRepository) FilterByFollowingId(followingId int) Relation {
	return Relation(func(entity string) string {
		var userIdField string
		if entity == "user" {
			userIdField = "id"
		} else {
			userIdField = "user_id"
		}
		return fmt.Sprintf(`
				INNER JOIN (
					SELECT follower_id FROM follows
					WHERE following_id=%d
				) ON follower_id = %ss.%s
			`,
			followingId,
			entity,
			userIdField,
		)
	})
}

func (r *followsRepository) WithTotalFollowers() Relation {
	return Relation(func(entity string) string {
		return fmt.Sprintf(`
				LEFT JOIN (
					SELECT following_id, COUNT(*) AS total_followers FROM follows GROUP BY following_id
				)
				ON %ss.id = following_id
			`,
			entity,
		)
	})
}

func (r *followsRepository) WithTotalFollowings() Relation {
	return Relation(func(entity string) string {
		return fmt.Sprintf(`
				LEFT JOIN (
					SELECT follower_id, COUNT(*) AS total_following FROM follows GROUP BY follower_id
				)
				ON %ss.id = follower_id
			`,
			entity,
		)
	})
}
