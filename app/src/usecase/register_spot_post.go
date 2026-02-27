package usecase

import (
	"context"
	"errors"
	"time"
	"app/domain/entities"
	"app/domain/value_objects"
)

type RegisterSpotPostInput struct {
	UserID    int
	SpotName  string
	Latitude  float64
	Longitude float64
	ImageURL  string
	Caption   string
	Overwrite bool
}

type RegisterSpotPostOutput struct {
	SpotID int
	PostID int
}

type RegisterSpotPostPresenter interface {
	Output(spot *entities.Spot, post *entities.Post) *RegisterSpotPostOutput
}

type RegisterSpotPostUseCase interface {
	Execute(ctx context.Context, input RegisterSpotPostInput) (*RegisterSpotPostOutput, error)
}

type registerSpotPostInteractor struct {
	presenter   RegisterSpotPostPresenter
	spotRepo    entities.SpotRepository
	postRepo    entities.PostRepository
}

func NewRegisterSpotPostInteractor(p RegisterSpotPostPresenter, s entities.SpotRepository, r entities.PostRepository) RegisterSpotPostUseCase {
	return &registerSpotPostInteractor{
		presenter: p,
		spotRepo:  s,
		postRepo:  r,
	}
}

func (i *registerSpotPostInteractor) Execute(ctx context.Context, input RegisterSpotPostInput) (*RegisterSpotPostOutput, error) {
	spot, err := entities.NewSpot(0, input.SpotName, input.Latitude, input.Longitude)
	if err != nil {
		return nil, err
	}
	// 既存スポット確認
	existing, err := i.spotRepo.FindByMeshID(spot.MeshID)
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 && !input.Overwrite {
		return nil, errors.New("spot already exists")
	}
	createdSpot, err := i.spotRepo.Create(spot)
	if err != nil {
		return nil, err
	}
	userID, err := entities.NewUser(input.UserID, "", "", "")
	if err != nil {
		return nil, err
	}
	imgURL, err := value_objects.NewImageURL(input.ImageURL)
	if err != nil {
		return nil, err
	}
	post, err := entities.NewPost(0, userID.ID.Value(), "", imgURL.String(), input.Caption, time.Now())
	if err != nil {
		return nil, err
	}
	createdPost, err := i.postRepo.Create(post)
	if err != nil {
		return nil, err
	}
	return i.presenter.Output(createdSpot, createdPost), nil
}
