package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lucaronca/wasa-homework/service/api/models"
	"github.com/lucaronca/wasa-homework/service/database"
)

type PhotosRepository interface {
	// Getters
	GetPhotoById(int) (*models.Photo, error)
	GetPhotos(int, int, ...Relation) (*[]models.Photo, error)
	GetPhotosCount(...Relation) (int, error)
	// Setters
	SetPhoto(string, int, time.Time) (int, error)
	RemovePhoto(int) error
	// Relation builders
	WithTotalPhotos() Relation
	FilterByPhotoId(int) Relation
}

type photosRepository struct {
	database.AppDatabase
}

func NewPhotosRepository(db database.AppDatabase) (PhotosRepository, error) {
	if db == nil {
		return nil, errors.New("database is required")
	}

	return &photosRepository{
		db,
	}, nil
}

func (r *photosRepository) GetPhotoById(photoId int) (*models.Photo, error) {
	var photo models.Photo
	var uploadDate string
	err := r.Conn().QueryRow(`
		SELECT photos.id, url, user_id, users.username, upload_date FROM photos
		INNER JOIN users ON users.id = user_id
		WHERE photos.id=?;
	`, photoId).Scan(&photo.Id, &photo.Url, &photo.Owner.Id, &photo.Owner.Username, &uploadDate)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	date, err := time.Parse(dateLayout, uploadDate)
	if err != nil {
		return nil, err
	}
	photo.UploadDate = date
	return &photo, nil
}

func (r *photosRepository) GetPhotos(offset, rowCount int, relations ...Relation) (*[]models.Photo, error) {
	q := queryBuilder("photo", relations...)
	rows, err := r.Conn().Query(fmt.Sprintf(`
		SELECT
			photos.id,
			url,
			user_id,
			users.username,
			upload_date,
			CASE WHEN total_likes IS NULL THEN 0 ELSE total_likes END as total_likes,
			CASE WHEN total_comments IS NULL THEN 0 ELSE total_comments END as total_comments,
			CASE WHEN user_liked_photo_id IS NULL THEN 0 ELSE 1 END as user_liked_photo
		FROM photos
		%s
		ORDER BY upload_date DESC
		LIMIT ? OFFSET ?;
	`, q), rowCount, offset)
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

	var photos []models.Photo
	for rows.Next() {
		var id int
		var url string
		var ownerId int
		var ownerUsername string
		var uploadDate string
		var totalLikes int
		var totalComments int
		var userLiked bool
		err = rows.Scan(&id, &url, &ownerId, &ownerUsername, &uploadDate, &totalLikes, &totalComments, &userLiked)
		if err != nil {
			return nil, err
		}

		date, err := time.Parse(dateLayout, uploadDate)
		if err != nil {
			return nil, err
		}

		photos = append(photos, models.Photo{
			Id:         id,
			Url:        url,
			UploadDate: date,
			Owner: models.BaseUser{
				Id:       ownerId,
				Username: ownerUsername,
			},
			TotalLikes:    totalLikes,
			TotalComments: totalComments,
			UserLiked:     userLiked,
		})
	}

	return &photos, nil
}

func (r *photosRepository) GetPhotosCount(relations ...Relation) (int, error) {
	q := queryBuilder("photo", relations...)
	var count int
	err := r.Conn().QueryRow(fmt.Sprintf(`
		SELECT COUNT(*) FROM photos
		%s;
	`, q)).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *photosRepository) SetPhoto(url string, userId int, time time.Time) (int, error) {
	result, err := r.Conn().Exec(`
		INSERT INTO photos (url, user_id, upload_date)
		VALUES (?, ?, ?);
	`, url, userId, time.Format(dateLayout))
	if err != nil {
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastInsertId), nil
}

func (r *photosRepository) RemovePhoto(photoId int) error {
	if _, err := r.Conn().Exec(`
		DELETE FROM photos
		WHERE id=?;
	`, photoId); err != nil {
		return err
	}
	return nil
}

func (r *photosRepository) WithTotalPhotos() Relation {
	return Relation(func(entity string) string {
		return fmt.Sprintf(`
				LEFT JOIN (
					SELECT user_id, COUNT(*) AS total_photos FROM photos GROUP BY user_id
				)
				ON %ss.id = user_id
			`,
			entity,
		)
	})
}

func (r *photosRepository) FilterByPhotoId(photoId int) Relation {
	return Relation(func(entity string) string {
		if entity == "photo" {
			return fmt.Sprintf(
				"WHERE id=%d",
				photoId,
			)
		}
		return fmt.Sprintf(
			"WHERE %ss.photo_id=%d",
			entity,
			photoId,
		)
	})
}
