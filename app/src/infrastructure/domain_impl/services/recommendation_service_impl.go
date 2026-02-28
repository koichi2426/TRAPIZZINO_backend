package domain_impl_services

import (
	"context"
	"fmt"
	"math"
	"app/domain/entities"
	"app/domain/services"
	"app/domain/value_objects"
)

type RecommendationServiceImpl struct {
	spotRepo entities.SpotRepository
}

func NewRecommendationServiceImpl(spotRepo entities.SpotRepository) services.RecommendationService {
	return &RecommendationServiceImpl{
		spotRepo: spotRepo,
	}
}

// Distill はドメインサービスのインターフェースに従い、6つのドメインオブジェクトとエラーを返します。
func (s *RecommendationServiceImpl) Distill(
	ctx context.Context,
	user *entities.User,
	lat value_objects.Latitude,
	lng value_objects.Longitude,
) (
	*entities.Spot,
	value_objects.TotalScore,
	value_objects.ResonanceCount,
	value_objects.DensityScore,
	value_objects.Reason,
	[]*entities.Post,
	error,
) {
	// ゼロ値の準備（エラー時用）
	var (
		emptySpot    *entities.Spot
		emptyScore   value_objects.TotalScore
		emptyRes     value_objects.ResonanceCount
		emptyDen     value_objects.DensityScore
		emptyReason  value_objects.Reason
		emptyPosts   []*entities.Post
	)

	// --- STEP 1: 空間の量子化 ---
	currentMesh, err := value_objects.NewMeshID(lat.Value(), lng.Value())
	if err != nil {
		return emptySpot, emptyScore, emptyRes, emptyDen, emptyReason, emptyPosts, err
	}

	// --- STEP 3: 共鳴者の特定 ---
	resonantUsers, err := s.spotRepo.FindResonantUsersWithMatchCount(ctx, user.ID)
	if err != nil {
		return emptySpot, emptyScore, emptyRes, emptyDen, emptyReason, emptyPosts, err
	}
	if len(resonantUsers) == 0 {
		return emptySpot, emptyScore, emptyRes, emptyDen, emptyReason, emptyPosts, fmt.Errorf("no resonant users found")
	}

	// --- STEP 4: 共鳴による「メッシュ代表店」の選定 ---
	targetMeshes := []value_objects.MeshID{currentMesh}
	
	resonantIDs := make([]value_objects.ID, 0, len(resonantUsers))
	resonanceMap := make(map[int]int)
	for _, ru := range resonantUsers {
		resonantIDs = append(resonantIDs, ru.ID)
		resonanceMap[ru.ID.Value()] = ru.MatchCount
	}

	allCandidateSpots, err := s.spotRepo.FindSpotsByMeshAndUsers(ctx, targetMeshes, resonantIDs)
	if err != nil {
		return emptySpot, emptyScore, emptyRes, emptyDen, emptyReason, emptyPosts, err
	}

	meshRepresentatives := make(map[string]*entities.Spot)
	meshTopResonance := make(map[string]int)

	for _, spot := range allCandidateSpots {
		mID := spot.MeshID.String()
		rCount := resonanceMap[spot.RegisteredUserID.Value()] 

		if rCount > meshTopResonance[mID] {
			meshTopResonance[mID] = rCount
			meshRepresentatives[mID] = spot
		}
	}

	// --- STEP 5 & 6: 統合スコア算出と運命の1軒の決定 ---
	var bestSpot *entities.Spot
	var maxScore float64
	var bestResonance int
	var bestDensity int

	for mID, spot := range meshRepresentatives {
		resCount := meshTopResonance[mID]
		density, _ := s.spotRepo.GetDensityScoreByMesh(ctx, spot.MeshID)

		dist := s.calculateDistance(lat.Value(), lng.Value(), spot.Latitude.Value(), spot.Longitude.Value())
		distanceWeight := 1.0 / (1.0 + math.Log1p(dist)) 

		scoreValue := (float64(resCount) * float64(density.Int())) * distanceWeight

		if scoreValue > maxScore {
			maxScore = scoreValue
			bestSpot = spot
			bestResonance = resCount
			bestDensity = density.Int()
		}
	}

	if bestSpot == nil {
		return emptySpot, emptyScore, emptyRes, emptyDen, emptyReason, emptyPosts, fmt.Errorf("could not distill the best spot")
	}

	// --- 戻り値の Value Object 化 ---
	totalScore, _ := value_objects.NewTotalScore(maxScore)
	resCountVO, _ := value_objects.NewResonanceCount(bestResonance)
	denScoreVO, _ := value_objects.NewDensityScore(bestDensity)
	reasonVO, _ := value_objects.NewReason(fmt.Sprintf(
		"あなたと%d地点で感性が一致する共鳴者が、激戦区（密度:%d）において選び抜いた至高の1軒です。",
		bestResonance,
		bestDensity,
	))

	// 関連する投稿（共鳴者の声）を取得
	posts, _ := s.spotRepo.FindPostsBySpot(ctx, bestSpot.ID)

	return bestSpot, totalScore, resCountVO, denScoreVO, reasonVO, posts, nil
}

func (s *RecommendationServiceImpl) calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371
	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}