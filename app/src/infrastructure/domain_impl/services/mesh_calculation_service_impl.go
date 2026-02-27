package domain_impl_services

import (
	"app/domain/services"
	"app/domain/value_objects"
)

// MeshCalculationServiceImplはMeshCalculationServiceインターフェースの具象実装です。
type MeshCalculationServiceImpl struct{}

func NewMeshCalculationServiceImpl() services.MeshCalculationService {
	return &MeshCalculationServiceImpl{}
}

// CalculateMeshIDは緯度経度からMeshIDを算出します（ダミー実装）。
func (s *MeshCalculationServiceImpl) CalculateMeshID(lat value_objects.Latitude, lng value_objects.Longitude) (value_objects.MeshID, error) {
	return value_objects.MeshID("dummy-mesh-id"), nil
}
