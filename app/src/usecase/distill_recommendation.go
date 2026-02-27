package usecase

import (
	"context"
	"app/domain/entities"
	"app/domain/services"
	"app/domain/value_objects"
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
	userID, err := value_objects.NewID(input.UserID)
	if err != nil {
		return nil, err
	}
	spot, err := i.spotRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	user := &entities.User{ID: userID}
	var spots []*entities.Spot
	if spot != nil {
		spots = append(spots, spot)
	}
	recommended, err := i.recommendation.RecommendSpot(user, spots)
	if err != nil {
		return nil, err
	}
	return i.presenter.Output(recommended), nil
}
