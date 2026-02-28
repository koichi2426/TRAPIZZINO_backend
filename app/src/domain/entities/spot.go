package entities

import (
	"context"
	"app/domain/value_objects"
)

type Spot struct {
	ID               value_objects.ID
	Name             value_objects.SpotName
	MeshID           value_objects.MeshID
	Latitude         value_objects.Latitude
	Longitude        value_objects.Longitude
	RegisteredUserID value_objects.ID 
}

func NewSpot(id int, name string, lat, lng float64, userID int) (*Spot, error) {
	spotID, err := value_objects.NewID(id)
	if err != nil {
		return nil, err
	}
	spotName, err := value_objects.NewSpotName(name)
	if err != nil {
		return nil, err
	}
	latitude, err := value_objects.NewLatitude(lat)
	if err != nil {
		return nil, err
	}
	longitude, err := value_objects.NewLongitude(lng)
	if err != nil {
		return nil, err
	}
	uID, err := value_objects.NewID(userID)
	if err != nil {
		return nil, err
	}
	meshID, err := value_objects.NewMeshID(lat, lng)
	if err != nil {
		return nil, err
	}

	return &Spot{
		ID:               spotID,
		Name:             spotName,
		MeshID:           meshID,
		Latitude:         latitude,
		Longitude:        longitude,
		RegisteredUserID: uID,
	}, nil
}

type ResonantUser struct {
	ID         value_objects.ID
	MatchCount int
}

type SpotRepository interface {
	Create(spot *Spot) (*Spot, error)
	// FindByID に ctx が必要なら追加（実装側が ScanContext 等を使う場合）
	FindByID(ctx context.Context, id value_objects.ID) (*Spot, error)
	FindByMeshID(meshID value_objects.MeshID) ([]*Spot, error)
	
	// 【修正】ビルドエラーの want 形式に合わせ、ctx を除外
	Update(spot *Spot) error
	Delete(id value_objects.ID) error

	// --- 蒸留アルゴリズム用メソッド ---
	FindResonantUsersWithMatchCount(ctx context.Context, userID value_objects.ID) ([]ResonantUser, error)
	FindSpotsByMeshAndUsers(ctx context.Context, meshIDs []value_objects.MeshID, userIDs []value_objects.ID) ([]*Spot, error)
	GetDensityScoreByMesh(ctx context.Context, meshID value_objects.MeshID) (value_objects.DensityScore, error)
	FindPostsBySpot(ctx context.Context, spotID value_objects.ID) ([]*Post, error)
}