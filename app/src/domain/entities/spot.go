package entities

import (
	"app/domain/value_objects"
)

type Spot struct {
	ID        value_objects.ID
	Name      value_objects.SpotName
	MeshID    value_objects.MeshID
	Latitude  value_objects.Latitude
	Longitude value_objects.Longitude
}

func NewSpot(id int, name string, lat, lng float64) (*Spot, error) {
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
	meshID, err := value_objects.NewMeshID(lat, lng)
	if err != nil {
		return nil, err
	}
	return &Spot{
		ID:        spotID,
		Name:      spotName,
		MeshID:    meshID,
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}

type SpotRepository interface {
	Create(spot *Spot) (*Spot, error)
	FindByID(id value_objects.ID) (*Spot, error)
	FindByMeshID(meshID value_objects.MeshID) ([]*Spot, error)
	Update(spot *Spot) error
	Delete(id value_objects.ID) error
}
