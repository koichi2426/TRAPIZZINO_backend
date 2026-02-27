package postgres

import (
	"context"
	"database/sql"
	"src/domain/entities"
	"src/domain/value_objects"
)

// UserRepositoryはPostgreSQLを用いたUserRepositoryの実装です。
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entities.User) (*entities.User, error) {
	query := `INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := r.db.QueryRow(query, user.Username, user.Email.String(), user.HashedPassword).Scan(&id)
	if err != nil {
		return nil, err
	}
	user.ID, _ = value_objects.NewID(id)
	return user, nil
}

func (r *UserRepository) FindByID(id value_objects.ID) (*entities.User, error) {
	query := `SELECT id, username, email, hashed_password FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id.Value())
	var uid int
	var username, email, hashedPassword string
	if err := row.Scan(&uid, &username, &email, &hashedPassword); err != nil {
		return nil, err
	}
	userID, _ := value_objects.NewID(uid)
	emailVO, _ := value_objects.NewEmail(email)
	return &entities.User{
		ID:             userID,
		Username:       username,
		Email:          emailVO,
		HashedPassword: hashedPassword,
	}, nil
}

func (r *UserRepository) FindByEmail(email value_objects.Email) (*entities.User, error) {
	query := `SELECT id, username, email, hashed_password FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email.String())
	var uid int
	var username, emailStr, hashedPassword string
	if err := row.Scan(&uid, &username, &emailStr, &hashedPassword); err != nil {
		return nil, err
	}
	userID, _ := value_objects.NewID(uid)
	emailVO, _ := value_objects.NewEmail(emailStr)
	return &entities.User{
		ID:             userID,
		Username:       username,
		Email:          emailVO,
		HashedPassword: hashedPassword,
	}, nil
}

func (r *UserRepository) Update(user *entities.User) error {
	query := `UPDATE users SET username = $1, email = $2, hashed_password = $3 WHERE id = $4`
	_, err := r.db.Exec(query, user.Username, user.Email.String(), user.HashedPassword, user.ID.Value())
	return err
}

func (r *UserRepository) Delete(id value_objects.ID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id.Value())
	return err
}
