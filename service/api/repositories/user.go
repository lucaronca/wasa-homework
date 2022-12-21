package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lucaronca/wasa-homework/service/api/models"
	"github.com/lucaronca/wasa-homework/service/database"
)

type UsersRepository interface {
	// Getters
	GetUserById(id int) (*models.BaseUser, error)
	GetUser(relations ...Relation) (*models.BaseUser, error)
	GetFullUser(relations ...Relation) (*models.FullUser, error)
	GetUsers(relations ...Relation) (*[]models.BaseUser, error)
	// Setters
	CreateUser(user *models.BaseUser) (userId int, err error)
	UpdateUser(user *models.BaseUser) (err error)
	// Relations builders
	WithUsers() Relation
	FilterByUserId(userId int) Relation
	FilterByUsername(username string, strict bool) Relation
}

type usersRepository struct {
	database.AppDatabase
}

func NewUsersRepository(db database.AppDatabase) (UsersRepository, error) {
	if db == nil {
		return nil, errors.New("database is required")
	}

	return &usersRepository{
		db,
	}, nil
}

func (r *usersRepository) GetUserById(id int) (*models.BaseUser, error) {
	var user models.BaseUser
	err := r.Conn().QueryRow(`
		SELECT id, username FROM users
		WHERE id=?
	`, id).Scan(&user.Id, &user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepository) GetFullUser(relations ...Relation) (*models.FullUser, error) {
	q := queryBuilder("user", relations...)
	var user models.FullUser
	err := r.Conn().QueryRow(fmt.Sprintf(`
		SELECT
			id,
			username,
			CASE WHEN total_followers IS NULL THEN 0 ELSE total_followers END as total_followers,
			CASE WHEN total_following IS NULL THEN 0 ELSE total_following END as total_followings,
			CASE WHEN total_photos IS NULL THEN 0 ELSE total_photos END as total_photos
		FROM users
		%s
		LIMIT 1;
	`, q)).Scan(
		&user.Id,
		&user.Username,
		&user.TotalFollowers,
		&user.TotalFollowings,
		&user.TotalPhotos,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepository) GetUserByUsername(username string) (*models.BaseUser, error) {
	var user models.BaseUser
	err := r.Conn().QueryRow(`
		SELECT id, username FROM users
		WHERE username=?
	`, username).Scan(&user.Id, &user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepository) CreateUser(user *models.BaseUser) (int, error) {
	result, err := r.Conn().Exec("INSERT INTO users (username) VALUES (?)", user.Username)
	if err != nil {
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastInsertId), nil
}

func (r *usersRepository) UpdateUser(user *models.BaseUser) error {
	_, err := r.Conn().Exec(`
		UPDATE users SET username=? WHERE id =?
	`, user.Username, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *usersRepository) GetUser(relations ...Relation) (*models.BaseUser, error) {
	q := queryBuilder("user", relations...)
	var user models.BaseUser
	err := r.Conn().QueryRow(fmt.Sprintf(`
		SELECT users.id, users.username FROM users
		%s
	`, q)).Scan(&user.Id, &user.Username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *usersRepository) GetUsers(relations ...Relation) (*[]models.BaseUser, error) {
	q := queryBuilder("user", relations...)
	rows, err := r.Conn().Query(
		fmt.Sprintf(
			"SELECT id, username FROM users %s",
			q,
		),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	var users []models.BaseUser
	for rows.Next() {
		var id int
		var username string
		err = rows.Scan(&id, &username)
		if err != nil {
			return nil, err
		}

		users = append(users, models.BaseUser{
			Id:       id,
			Username: username,
		})
	}

	return &users, nil
}

// Relations builders
func (r *usersRepository) WithUsers() Relation {
	return Relation(func(entity string) string {
		return fmt.Sprintf(
			"INNER JOIN users ON users.id = %ss.user_id",
			entity,
		)
	})
}

func (r *usersRepository) FilterByUserId(userId int) Relation {
	return Relation(func(entity string) string {
		if entity == "user" {
			return fmt.Sprintf(
				"WHERE id=%d",
				userId,
			)
		}
		return fmt.Sprintf(
			"WHERE %ss.user_id=%d",
			entity,
			userId,
		)
	})
}

func (r *usersRepository) FilterByUsername(username string, strict bool) Relation {
	return Relation(func(entity string) string {
		if strict {
			return fmt.Sprintf(
				"WHERE username=\"%s\"",
				username,
			)
		}
		return "WHERE username LIKE " + "\"%" + username + "%\";"
	})
}
