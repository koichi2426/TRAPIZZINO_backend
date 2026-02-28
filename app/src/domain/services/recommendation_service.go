package services

import (
	"context"
	"app/domain/entities"
	"app/domain/value_objects"
)

type RecommendationService interface {
	// Distill は、蒸留メッシュアルゴリズムを用いて「運命の1軒」を算出します。
	// 循環参照を回避し、かつドメイン層の純粋性を保つため、
	// 構造体（DTO）を介さず各ドメインオブジェクトを個別に返却します。
	Distill(
		ctx context.Context, 
		user *entities.User, 
		lat value_objects.Latitude, 
		lng value_objects.Longitude,
	) (
		spot *entities.Spot,
		totalScore value_objects.TotalScore,
		resonanceCount value_objects.ResonanceCount,
		densityScore value_objects.DensityScore,
		reason value_objects.Reason,
		posts []*entities.Post,
		err error,
	)
}