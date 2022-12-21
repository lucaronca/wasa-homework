package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lucaronca/wasa-homework/service/api/models"
	"github.com/lucaronca/wasa-homework/service/database"
)

type CommentsRepository interface {
	// Getters
	GetCommentById(int, ...Relation) (*models.Comment, error)
	GetComments(relations ...Relation) (*[]models.Comment, error)
	// Setters
	SetComment(int, int, time.Time, string) (int, error)
	RemoveComment(int) error
	// Relation builders
	WithTotalComments() Relation
}

type commentsRepository struct {
	database.AppDatabase
}

func NewCommentsRepository(db database.AppDatabase) (CommentsRepository, error) {
	if db == nil {
		return nil, errors.New("database is required")
	}

	return &commentsRepository{
		db,
	}, nil
}

func (r *commentsRepository) GetCommentById(id int, relations ...Relation) (*models.Comment, error) {
	q := queryBuilder("comment", relations...)
	var comment models.Comment
	var date string
	err := r.Conn().QueryRow(fmt.Sprintf(`
		SELECT comments.id, photo_id, user_id, username, date, content FROM comments
		%s
		WHERE comments.id=?
	`, q), id).Scan(&comment.Id, &comment.Photo.Id, &comment.Owner.Id, &comment.Owner.Username, &date, &comment.Content)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	parsedDate, err := time.Parse(dateLayout, date)
	if err != nil {
		return nil, err
	}

	comment.Date = parsedDate

	return &comment, nil
}

func (r *commentsRepository) SetComment(photoId int, userId int, time time.Time, content string) (int, error) {
	result, err := r.Conn().Exec(`
		INSERT INTO comments (photo_id, user_id, date, content)
		VALUES (?, ?, ?, ?)
	`, photoId, userId, time.Format(dateLayout), content)
	if err != nil {
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastInsertId), nil
}

func (r *commentsRepository) RemoveComment(commentId int) error {
	if _, err := r.Conn().Exec(`
		DELETE FROM comments
		WHERE id=?;
	`, commentId); err != nil {
		return err
	}
	return nil
}

func (r *commentsRepository) GetComments(relations ...Relation) (*[]models.Comment, error) {
	q := queryBuilder("comment", relations...)
	rows, err := r.Conn().Query(fmt.Sprintf(`
		SELECT
			comments.id,
			comments.photo_id,
			comments.user_id,
			username,
			comments.date,
			comments.content
		FROM comments
		%s
		ORDER BY comments.date DESC
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

	var comments []models.Comment
	for rows.Next() {
		var id int
		var photoId int
		var ownerId int
		var ownerUsername string
		var date string
		var content string
		err = rows.Scan(&id, &photoId, &ownerId, &ownerUsername, &date, &content)
		if err != nil {
			return nil, err
		}

		parsedDate, err := time.Parse(dateLayout, date)
		if err != nil {
			return nil, err
		}

		comments = append(comments, models.Comment{
			Id:      id,
			Date:    parsedDate,
			Content: content,
			Photo: models.BasePhoto{
				Id: photoId,
			},
			Owner: models.BaseUser{
				Id:       ownerId,
				Username: ownerUsername,
			},
		})
	}

	return &comments, nil
}

// Relations builders
func (r *commentsRepository) WithTotalComments() Relation {
	return Relation(func(entity string) string {
		return fmt.Sprintf(`
				LEFT JOIN (
					SELECT %v_id as comments_%[1]vs_id, COUNT(*) AS total_comments FROM comments GROUP BY %[1]v_id
				)
				ON %[1]vs.id = comments_%[1]vs_id
			`,
			entity,
		)
	})
}
