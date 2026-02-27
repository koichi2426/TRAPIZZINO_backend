package postgres

import (
	"database/sql"
	"app/domain/entities"
	"app/domain/value_objects"
)

// PostRepositoryはPostgreSQLを用いたPostRepositoryの実装です。
type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *entities.Post) (*entities.Post, error) {
	query := `INSERT INTO posts (user_id, spot_id, image_url, caption, posted_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int
	err := r.db.QueryRow(query, post.UserID.Value(), post.ID.Value(), post.ImageURL.String(), post.Caption, post.PostedAt).Scan(&id)
	if err != nil {
		return nil, err
	}
	post.ID, _ = value_objects.NewID(id)
	return post, nil
}

func (r *PostRepository) FindByID(id value_objects.ID) (*entities.Post, error) {
	query := `SELECT id, user_id, spot_id, image_url, caption, posted_at FROM posts WHERE id = $1`
	row := r.db.QueryRow(query, id.Value())
	var pid, userID, spotID int
	var imageURL, caption string
	var postedAt sql.NullTime
	if err := row.Scan(&pid, &userID, &spotID, &imageURL, &caption, &postedAt); err != nil {
		return nil, err
	}
	postID, _ := value_objects.NewID(pid)
	userIDVO, _ := value_objects.NewID(userID)
	imgURL, _ := value_objects.NewImageURL(imageURL)
	capVO, _ := value_objects.NewCaption(caption)
	return &entities.Post{
		ID:       postID,
		UserID:   userIDVO,
		ImageURL: imgURL,
		Caption:  capVO,
		PostedAt: postedAt.Time,
	}, nil
}

func (r *PostRepository) FindBySpotID(spotID value_objects.ID) ([]*entities.Post, error) {
	query := `SELECT id, user_id, spot_id, image_url, caption, posted_at FROM posts WHERE spot_id = $1`
	rows, err := r.db.Query(query, spotID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*entities.Post
	for rows.Next() {
		var pid, userID, spotID int
		var imageURL, caption string
		var postedAt sql.NullTime
		if err := rows.Scan(&pid, &userID, &spotID, &imageURL, &caption, &postedAt); err != nil {
			return nil, err
		}
		postID, _ := value_objects.NewID(pid)
		userIDVO, _ := value_objects.NewID(userID)
		imgURL, _ := value_objects.NewImageURL(imageURL)
		capVO, _ := value_objects.NewCaption(caption)
		posts = append(posts, &entities.Post{
			ID:       postID,
			UserID:   userIDVO,
			ImageURL: imgURL,
			Caption:  capVO,
			PostedAt: postedAt.Time,
		})
	}
	return posts, nil
}

func (r *PostRepository) Update(post *entities.Post) error {
	query := `UPDATE posts SET image_url = $1, caption = $2, posted_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, post.ImageURL.String(), post.Caption, post.PostedAt, post.ID.Value())
	return err
}

func (r *PostRepository) Delete(id value_objects.ID) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.Exec(query, id.Value())
	return err
}
