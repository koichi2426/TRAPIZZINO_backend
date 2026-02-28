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
    latKey := math.Floor(lat * 100)
    lngKey := math.Floor(lng * 100)

    // 正の値になるようオフセットを加え、ユニークなID文字列を作成
    mesh := fmt.Sprintf("MSH-%05.0f-%05.0f", latKey+9000, lngKey+18000)
    
    return MeshID(mesh), nil
}

// GetSurroundingMeshIDs は現在のメッシュに隣接する 8 方向のメッシュ ID を返します
func (m MeshID) GetSurroundingMeshIDs() []MeshID {
    var latKey, lngKey float64
    // 現在の文字列 ID から座標キーを抽出
    _, err := fmt.Sscanf(string(m), "MSH-%f-%f", &latKey, &lngKey)
    if err != nil {
        return nil
    }

    surroundings := make([]MeshID, 0, 8)

    // 自身の周囲 1 マス分（計 9 マス）をループし、自分以外を追加
    for dLat := -1.0; dLat <= 1.0; dLat++ {
        for dLng := -1.0; dLng <= 1.0; dLng++ {
            if dLat == 0 && dLng == 0 {
                continue
            }
            // 隣接する座標キーで ID を再構成
            mesh := fmt.Sprintf("MSH-%05.0f-%05.0f", latKey+dLat, lngKey+dLng)
            surroundings = append(surroundings, MeshID(mesh))
        }
    }
    return surroundings
}

func (m MeshID) String() string {
    return string(m)
}