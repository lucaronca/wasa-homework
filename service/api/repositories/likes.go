package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lucaronca/wasa-homework/service/api/models"
	"github.com/lucaronca/wasa-homework/service/database"
)

type LikesRepository interface {
	// Getters
	GetLikes(relations ...Relation) (*[]models.Like, error)
	// Setters
	SetLike(int, int, time.Time) error
	RemoveLike(int, int) error
	// Relation builders
	WithTotalLikes() Relation
}

type likesRepository struct {
	database.AppDatabase
}

func NewLikesRepository(db database.AppDatabase) (LikesRepository, error) {
	if db == nil {
		return nil, errors.New("database is required")
	}

	return &likesRepository{
		db,
	}, nil
}

func (r *likesRepository) SetLike(photoId int, userId int, time time.Time) error {
	if _, err := r.Conn().Exec(`
		INSERT OR IGNORE INTO likes (photo_id, user_id, date)
		VALUES (?, ?, ?)
	`, photoId, userId, time.Format(dateLayout)); err != nil {
		return err
	}
	return nil
}

func (r *likesRepository) RemoveLike(photoId int, userId int) error {
	if _, err := r.Conn().Exec(`
		DELETE FROM likes
		WHERE photo_id=? AND user_id=?;
	`, photoId, userId); err != nil {
		return err
	}
	return nil
}

func (r *likesRepository) GetLikes(relations ...Relation) (*[]models.Like, error) {
	q := queryBuilder("like", relations...)
	rows, err := r.Conn().Query(fmt.Sprintf(`
		SELECT likes.id, likes.photo_id, likes.user_id, username, likes.date FROM likes
		%s
		ORDER BY date DESC
	`, q))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	var likes []models.Like
	for rows.Next() {
		var id int
		var photoId int
		var ownerId int
		var ownerUsername string
		var date string
		err = rows.Scan(&id, &photoId, &ownerId, &ownerUsername, &date)
		if err != nil {
			return nil, err
		}

		parsedDate, err := time.Parse(dateLayout, date)
		if err != nil {
			return nil, err
		}

		likes = append(likes, models.Like{
			Id:   id,
			Date: parsedDate,
			Photo: models.BasePhoto{
				Id: photoId,
			},
			Owner: models.BaseUser{
				Id:       ownerId,
				Username: ownerUsername,
			},
		})
	}

	return &likes, nil
}

// Relations builders
func (r *likesRepository) WithTotalLikes() Relation {
	return Relation(func(entity string) string {
		return fmt.Sprintf(`
				LEFT JOIN (
					SELECT %v_id as likes_%[1]vs_id, COUNT(*) AS total_likes FROM likes GROUP BY %[1]v_id
				)
				ON %[1]vs.id = likes_%[1]vs_id
			`,
			entity,
		)
	})
}
