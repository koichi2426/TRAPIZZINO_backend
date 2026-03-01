package usecase

import (
	"context"
	"app/domain/entities"
	"app/domain/services"
	"app/domain/value_objects"
)

type ListMySpotsInput struct {
	Token string
}

type SpotPostPair struct {
	Spot *entities.Spot
	Post *entities.Post
}

type ListMySpotsOutput struct {
	Pairs []SpotPostPair
}

type ListMySpotsPresenter interface {
	Output(pairs []SpotPostPair) *ListMySpotsOutput
}

type ListMySpotsUseCase interface {
	Execute(ctx context.Context, input ListMySpotsInput) (*ListMySpotsOutput, error)
}

type listMySpotsInteractor struct {
	presenter ListMySpotsPresenter
	spotRepo  entities.SpotRepository
	postRepo  entities.PostRepository
	authService services.AuthDomainService
}

func NewListMySpotsInteractor(
	p ListMySpotsPresenter, 
	s entities.SpotRepository, 
	r entities.PostRepository,
	a services.AuthDomainService,
) ListMySpotsUseCase {
	return &listMySpotsInteractor{
		presenter:   p,
		spotRepo:    s,
		postRepo:    r,
		authService: a,
	}
}

func (i *listMySpotsInteractor) Execute(ctx context.Context, input ListMySpotsInput) (*ListMySpotsOutput, error) {
	user, err := i.authService.VerifyToken(ctx, input.Token)
	if err != nil {
		return nil, err
	}

	userID, err := value_objects.NewID(user.ID.Value())
	if err != nil {
		return nil, err
	}

	// 修正ポイント：インターフェースの変更に合わせて第一引数に ctx を追加
	spot, err := i.spotRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var pairs []SpotPostPair
	if spot != nil {
		posts, err := i.postRepo.FindBySpotID(spot.ID)
		if err != nil {
			return nil, err
		}
		for _, post := range posts {
			pairs = append(pairs, SpotPostPair{Spot: spot, Post: post})
		}
	}
	return i.presenter.Output(pairs), nil
}