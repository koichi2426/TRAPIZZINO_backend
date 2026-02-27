package domain_impl_services

import (
	"app/domain/entities"
	"app/domain/services"
)

// RecommendationServiceImplはRecommendationServiceインターフェースの具象実装です。
type RecommendationServiceImpl struct{}

func NewRecommendationServiceImpl() services.RecommendationService {
	return &RecommendationServiceImpl{}
}

// RecommendSpotは最適なスポットを返します（ダミー実装）。
func (s *RecommendationServiceImpl) RecommendSpot(user *entities.User, spots []*entities.Spot) (*entities.Spot, error) {
	if len(spots) == 0 {
		return nil, nil
	}
	return spots[0], nil
}
