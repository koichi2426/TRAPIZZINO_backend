package usecase

import (
	"context"
	"src/domain/entities"
	"src/domain/services"
)

type DistillRecommendationInput struct {
	UserID int
}

type DistillRecommendationOutput struct {
	SpotID   int
	SpotName string
}

type DistillRecommendationPresenter interface {
	Output(spot *entities.Spot) *DistillRecommendationOutput
}

type DistillRecommendationUseCase interface {
	Execute(ctx context.Context, input DistillRecommendationInput) (*DistillRecommendationOutput, error)
}

type distillRecommendationInteractor struct {
	presenter      DistillRecommendationPresenter
	recommendation services.RecommendationService
	spotRepo       entities.SpotRepository
}

func NewDistillRecommendationInteractor(p DistillRecommendationPresenter, r services.RecommendationService, s entities.SpotRepository) DistillRecommendationUseCase {
	return &distillRecommendationInteractor{
		presenter:      p,
		recommendation: r,
		spotRepo:       s,
	}
}

func (i *distillRecommendationInteractor) Execute(ctx context.Context, input DistillRecommendationInput) (*DistillRecommendationOutput, error) {
	spots, err := i.spotRepo.FindByID(input.UserID)
	if err != nil {
		return nil, err
	}
	user := &entities.User{ID: value_objects.ID(input.UserID)}
	recommended, err := i.recommendation.RecommendSpot(user, spots)
	if err != nil {
		return nil, err
	}
	return i.presenter.Output(recommended), nil
}
