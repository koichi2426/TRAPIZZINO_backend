package domain_impl_services

import (
	"context"
	"fmt"
	"math"
	"app/src/domain/entities"
	"app/src/domain/services"
	"app/src/domain/value_objects"
)

type RecommendationServiceImpl struct {
	spotRepo entities.SpotRepository
}

func NewRecommendationServiceImpl(spotRepo entities.SpotRepository) services.RecommendationService {
	return &RecommendationServiceImpl{
		spotRepo: spotRepo,
	}
}

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

	// --- STEP 3: 共鳴者の特定（店舗一致に基づくメンター抽出） ---
	resonantUsers, err := s.spotRepo.FindResonantUsersWithMatchCount(ctx, user.ID)
	if err != nil {
		return emptySpot, emptyScore, emptyRes, emptyDen, emptyReason, emptyPosts, err
	}
	if len(resonantUsers) == 0 {
		return emptySpot, emptyScore, emptyRes, emptyDen, emptyReason, emptyPosts, fmt.Errorf("no resonant users found")
	}

	// --- STEP 4: 共鳴による「メッシュ代表店」の選定 ---
	targetMeshes := append([]value_objects.MeshID{currentMesh}, currentMesh.GetSurroundingMeshIDs()...)
	
	resonanceMap := make(map[int]int)
	resonantIDs := make([]value_objects.ID, 0, len(resonantUsers))
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

		// 店舗一致（resCount）の価値を最大化する重み付け
		resonanceWeight := (math.Log1p(float64(resCount)) * 3.0) + 1.0
		scoreValue := (resonanceWeight * float64(density.Int())) * distanceWeight

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

	// --- 最終結果のパッキングと推薦理由の生成 ---
	totalScore, _ := value_objects.NewTotalScore(maxScore)
	resCountVO, _ := value_objects.NewResonanceCount(bestResonance)
	denScoreVO, _ := value_objects.NewDensityScore(bestDensity)
	reasonVO, _ := value_objects.NewReason(fmt.Sprintf(
		"あなたと %d 箇所で『全く同じ一軒』を選び抜いた共鳴者が、激戦区（熱量:%d）で王座に据えた至高の1軒です。",
		bestResonance,
		bestDensity,
	))

	// --- 共鳴者による投稿内容の抽出 ---
	allPosts, _ := s.spotRepo.FindPostsBySpot(ctx, bestSpot.ID)
	var resonantPosts []*entities.Post
	for _, p := range allPosts {
		if _, ok := resonanceMap[p.UserID.Value()]; ok {
			resonantPosts = append(resonantPosts, p)
		}
	}

	return bestSpot, totalScore, resCountVO, denScoreVO, reasonVO, resonantPosts, nil
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