package services

import (
	"app/domain/entities"
)

type RecommendationService interface {
	RecommendSpot(user *entities.User, spots []*entities.Spot) (*entities.Spot, error)
}
