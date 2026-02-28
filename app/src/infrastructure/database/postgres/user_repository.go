package postgres

import (
	"context"
	"database/sql"
	"app/domain/entities"
	"app/domain/value_objects"
)

type UserRepository struct {
	db *sql.DB
}

// 戻り値をインターフェース型に合わせるのが一般的ですが、
// router.go の記述に合わせて一旦このままにします。
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entities.User) (*entities.User, error) {
	query := `INSERT INTO users (username, email, hashed_password) VALUES ($1, $2, $3) RETURNING id`
	var id int
	// 修正: Username と HashedPassword も .String() で生の文字列を渡す
	err := r.db.QueryRow(query, 
		user.Username.String(), 
		user.Email.String(), 
		user.HashedPassword.String(),
	).Scan(&id)
	
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
	uname, _ := value_objects.NewUsername(username)
	hashVO, _ := value_objects.NewHashedPassword(hashedPassword)
	return &entities.User{
		ID:             userID,
		Username:       uname,
		Email:          emailVO,
		HashedPassword: hashVO,
	}, nil
}

func (r *UserRepository) FindByEmail(email value_objects.Email) (*entities.User, error) {
	query := `SELECT id, username, email, hashed_password FROM users WHERE email = $1`
	// email.String() は OK
	row := r.db.QueryRow(query, email.String())
	var uid int
	var username, emailStr, hashedPassword string
	if err := row.Scan(&uid, &username, &emailStr, &hashedPassword); err != nil {
		return nil, err
	}
	userID, _ := value_objects.NewID(uid)
	emailVO, _ := value_objects.NewEmail(emailStr)
	uname, _ := value_objects.NewUsername(username)
	hashVO, _ := value_objects.NewHashedPassword(hashedPassword)
	return &entities.User{
		ID:             userID,
		Username:       uname,
		Email:          emailVO,
		HashedPassword: hashVO,
	}, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `SELECT id, username, email, hashed_password FROM users WHERE username = $1`
	row := r.db.QueryRowContext(ctx, query, username)
	
	var uid int
	var unameStr, emailStr, hashedPassword string
	if err := row.Scan(&uid, &unameStr, &emailStr, &hashedPassword); err != nil {
		return nil, err
	}

	userID, _ := value_objects.NewID(uid)
	emailVO, _ := value_objects.NewEmail(emailStr)
	uname, _ := value_objects.NewUsername(unameStr)
	hashVO, _ := value_objects.NewHashedPassword(hashedPassword)
	
	return &entities.User{
		ID:             userID,
		Username:       uname,
		Email:          emailVO,
		HashedPassword: hashVO,
	}, nil
}

func (r *UserRepository) Update(user *entities.User) error {
	query := `UPDATE users SET username = $1, email = $2, hashed_password = $3 WHERE id = $4`
	// 修正: 全て String() / Value() を介して渡す
	_, err := r.db.Exec(query, 
		user.Username.String(), 
		user.Email.String(), 
		user.HashedPassword.String(), 
		user.ID.Value(),
	)
	return err
}

func (r *UserRepository) Delete(id value_objects.ID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id.Value())
	return err
}