package value_objects

import (
	"errors"
	"fmt"
)

type MeshID string

// NewMeshID: メッシュIDの生成とバリデーション
func NewMeshID(lat, lng float64) (MeshID, error) {
	// 例: 1kmメッシュコード（簡易実装、実際の仕様に合わせて調整）
	if lat < -90 || lat > 90 {
		return "", errors.New("latitude out of range")
	}
	if lng < -180 || lng > 180 {
		return "", errors.New("longitude out of range")
	}
	// メッシュID生成ロジック（例: 小数点以下2桁で連結）
	mesh := fmt.Sprintf("%02d%03d", int(lat+90), int(lng+180))
	return MeshID(mesh), nil
}

func (m MeshID) String() string {
	return string(m)
}
