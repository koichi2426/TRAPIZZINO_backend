package usecase

import (
	"context"
	"src/domain/entities"
)

type ListMySpotsInput struct {
	UserID int
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
}

func NewListMySpotsInteractor(p ListMySpotsPresenter, s entities.SpotRepository, r entities.PostRepository) ListMySpotsUseCase {
	return &listMySpotsInteractor{
		presenter: p,
		spotRepo:  s,
		postRepo:  r,
	}
}

func (i *listMySpotsInteractor) Execute(ctx context.Context, input ListMySpotsInput) (*ListMySpotsOutput, error) {
	spots, err := i.spotRepo.FindByID(input.UserID)
	if err != nil {
		return nil, err
	}
	var pairs []SpotPostPair
	for _, spot := range spots {
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
