package services

import (
	"src/domain/entities"
)

type RecommendationService interface {
	RecommendSpot(user *entities.User, spots []*entities.Spot) (*entities.Spot, error)
}
