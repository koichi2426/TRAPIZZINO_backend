package entities

import (
	"time"
	"app/domain/value_objects"
)

type Post struct {
	ID        value_objects.ID
	UserID    value_objects.ID
	Username  value_objects.Username
	ImageURL  value_objects.ImageURL
	Caption   value_objects.Caption
	PostedAt  time.Time
}

func NewPost(id, userID int, username, imageURL, caption string, postedAt time.Time) (*Post, error) {
	pid, err := value_objects.NewID(id)
	if err != nil {
		return nil, err
	}
	uid, err := value_objects.NewID(userID)
	if err != nil {
		return nil, err
	}
	uname, err := value_objects.NewUsername(username)
	if err != nil {
		return nil, err
	}
	imgURL, err := value_objects.NewImageURL(imageURL)
	if err != nil {
		return nil, err
	}
	capVO, err := value_objects.NewCaption(caption)
	if err != nil {
		return nil, err
	}
	return &Post{
		ID:       pid,
		UserID:   uid,
		Username: uname,
		ImageURL: imgURL,
		Caption:  capVO,
		PostedAt: postedAt,
	}, nil
}

type PostRepository interface {
	Create(post *Post) (*Post, error)
	FindByID(id value_objects.ID) (*Post, error)
	FindBySpotID(spotID value_objects.ID) ([]*Post, error)
	Update(post *Post) error
	Delete(id value_objects.ID) error
}
