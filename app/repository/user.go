package repository

import (
	"database/sql"
	"fmt"
	"golang-test1/app/model"
	"golang-test1/platform/database"
	"time"

	"github.com/google/uuid"
)

type UserRepo struct {
	db *database.DB
}

func NewUserRepo(db *database.DB) UserRepository {
	return &UserRepo{db}
}

func (u UserRepo) NewUserRepo(db *database.DB) UserRepository {
	return &UserRepo{db}
}

func (u UserRepo) GetByUsername(username string) (*model.User, error) {
	user := model.User{}
	query := `SELECT * FROM "users" WHERE username = $1 AND is_deleted=FALSE`
	err := u.db.Get(&user, query, username)

	return &user, err
}

func (repo *UserRepo) Exists(username, email string) (bool, error) {
	findByUser := len(username)
	findByEmail := len(email)
	if findByUser <= 0 || findByEmail <= 0 {
		return false, nil
	}
	var placeHolderValues []interface{}
	query := `SELECT 1 FROM "users" WHERE is_deleted = FALSE`
	if findByUser > 0 && findByEmail > 0 {
		query = fmt.Sprintf("%s WHERE username = $1 OR email = $2", query)
		placeHolderValues = append(placeHolderValues, username, email)
	} else {
		if findByUser > 0 {
			query = fmt.Sprintf("%s WHERE username = $1", query)
			placeHolderValues = append(placeHolderValues, username)
		} else {
			query = fmt.Sprintf("%s WHERE email = $1", query)
			placeHolderValues = append(placeHolderValues, email)
		}
	}

	var exists bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := repo.db.QueryRow(query, placeHolderValues...).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}
func (repo *UserRepo) Create(u *model.RegisterUser) error {
	query := `INSERT INTO users (id, username, email, password_hash, full_name, role, is_active, created_at, updated_at, is_deleted)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := repo.db.Exec(query, uuid.New(), u.UserName, u.Email, u.Password, u.FullName, u.Role, true, time.Now(), time.Now(), false)
	return err
}
func (repo *UserRepo) All(limit int, offset uint) ([]*model.User, error) {
	var users []*model.User
	query := `SELECT * FROM "users" WHERE is_deleted=FALSE`
	var err error

	if limit > 0 && offset >= 0 {
		query = fmt.Sprintf("%s LIMIT $1 OFFSET $2", query)
		err = repo.db.Select(&users, query, limit, offset)
	} else {
		err = repo.db.Select(&users, query)
	}

	return users, err
}
func (repo *UserRepo) Get(ID string) (*model.User, error) {
	user := model.User{}
	query := `SELECT * FROM "users" WHERE id = $1 AND is_deleted = FALSE`
	err := repo.db.Get(&user, query, ID)

	return &user, err
}
func (repo *UserRepo) Update(ID string, u *model.UpdateUser) error {
	query := `UPDATE users SET 
        username = COALESCE($2, username),
        email = COALESCE($3, email),
        full_name = COALESCE($4, full_name), 
        role = COALESCE($5, role),
        updated_at = $6
        WHERE id = $1 AND is_deleted = FALSE`

	_, err := repo.db.Exec(
		query,
		ID,
		u.UserName,
		u.Email,
		u.FullName,
		u.Role,
		time.Now(),
	)
	return err
}

// The Delete function is correct, just needs to match the field names
func (repo *UserRepo) Delete(ID string) error {
	query := `UPDATE users SET is_deleted = TRUE, updated_at = $2 WHERE id = $1`
	_, err := repo.db.Exec(query, ID, time.Now())
	return err
}

// ChangePassword updates a user's password
func (repo *UserRepo) ChangePassword(ID string, newPassword string) error {
	query := `UPDATE users SET password_hash = $2, updated_at = $3 WHERE id = $1 AND is_deleted = FALSE`
	_, err := repo.db.Exec(query, ID, newPassword, time.Now())
	return err
}

// GetPassword retrieves a user's hashed password for verification
func (repo *UserRepo) GetPassword(ID string) (string, error) {
	var passwordHash string
	query := `SELECT password_hash FROM users WHERE id = $1 AND is_deleted = FALSE`
	err := repo.db.QueryRow(query, ID).Scan(&passwordHash)
	return passwordHash, err
}
