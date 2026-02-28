package value_objects

import (
	"errors"
	"fmt"
	"math"
)

type MeshID string

// NewMeshID は、緯度経度から「蒸留」の基準となる 1km メッシュ ID を生成します。
func NewMeshID(lat, lng float64) (MeshID, error) {
	// バリデーション
	if lat < -90 || lat > 90 {
		return "", errors.New("latitude out of range")
	}
	if lng < -180 || lng > 180 {
		return "", errors.New("longitude out of range")
	}

	// 0.01度単位（約1km四方）で座標を固定（空間の量子化）
	// math.Floor を使うことで、同じグリッド内なら必ず同じ値になります
	latKey := math.Floor(lat * 100)
	lngKey := math.Floor(lng * 100)

	// 正の値になるようオフセットを加え、ユニークなID文字列を作成
	mesh := fmt.Sprintf("MSH-%05.0f-%05.0f", latKey+9000, lngKey+18000)
	
	return MeshID(mesh), nil
}

func (m MeshID) String() string {
	return string(m)
}