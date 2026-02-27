package services

import "src/domain/value_objects"

type MeshCalculationService interface {
	CalculateMeshID(lat value_objects.Latitude, lng value_objects.Longitude) (value_objects.MeshID, error)
}
