package postgres

import (
	"app/domain/entities"
	"app/domain/value_objects"
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type spotRepository struct {
	db *sql.DB
}

func NewSpotRepository(db *sql.DB) entities.SpotRepository {
	return &spotRepository{db: db}
}

// 【修正】ST_DWithin(1m) ではなく、MeshID による厳密な区画判定に変更
func (r *spotRepository) FindByLocation(ctx context.Context, lat, lng float64) (*entities.Spot, error) {
	// 緯度経度から MeshID を生成
	meshVO, _ := value_objects.NewMeshID(lat, lng)

	query := `
        SELECT id, name, ST_X(location::geometry), ST_Y(location::geometry), registered_user_id 
        FROM spots 
        WHERE mesh_id = $1 
        LIMIT 1`

	var sid, uid int
	var name string
	var rLng, rLat float64

	err := r.db.QueryRowContext(ctx, query, meshVO.String()).Scan(&sid, &name, &rLng, &rLat, &uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return entities.NewSpot(sid, name, rLat, rLng, uid)
}

// 【修正】ON CONFLICT から registered_user_id を削除。MeshID のみで一意性を担保。
func (r *spotRepository) Create(spot *entities.Spot) (*entities.Spot, error) {
	query := `INSERT INTO spots (name, mesh_id, location, registered_user_id) 
    VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326), $5)
    ON CONFLICT (mesh_id) DO UPDATE 
    SET name = EXCLUDED.name, location = EXCLUDED.location 
    RETURNING id`

	var id int
	err := r.db.QueryRow(query,
		spot.Name.String(),
		spot.MeshID.String(),
		spot.Longitude.Value(),
		spot.Latitude.Value(),
		spot.RegisteredUserID.Value(),
	).Scan(&id)

	if err != nil {
		return nil, err
	}
	spot.ID, _ = value_objects.NewID(id)
	return spot, nil
}

// --- 以下、変更なしのメソッド群 ---

func (r *spotRepository) FindByID(ctx context.Context, id value_objects.ID) (*entities.Spot, error) {
	query := `SELECT id, name, ST_X(location::geometry), ST_Y(location::geometry), registered_user_id FROM spots WHERE id = $1`
	var sid, uid int
	var name string
	var lng, lat float64
	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(&sid, &name, &lng, &lat, &uid)
	if err != nil {
		return nil, err
	}
	return entities.NewSpot(sid, name, lat, lng, uid)
}

func (r *spotRepository) FindByMeshID(meshID value_objects.MeshID) ([]*entities.Spot, error) {
	query := `SELECT id, name, ST_X(location::geometry), ST_Y(location::geometry), registered_user_id FROM spots WHERE mesh_id = $1`
	rows, err := r.db.Query(query, meshID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spots []*entities.Spot
	for rows.Next() {
		var sid, uid int
		var name string
		var lng, lat float64
		if err := rows.Scan(&sid, &name, &lng, &lat, &uid); err != nil {
			return nil, err
		}
		s, _ := entities.NewSpot(sid, name, lat, lng, uid)
		spots = append(spots, s)
	}
	return spots, nil
}

func (r *spotRepository) FindResonantUsersWithMatchCount(ctx context.Context, userID value_objects.ID) ([]entities.ResonantUser, error) {
	query := `
        SELECT p.user_id, COUNT(DISTINCT s.mesh_id) as match_count 
        FROM posts p
        JOIN spots s ON p.spot_id = s.id
        WHERE s.mesh_id IN (
            SELECT s2.mesh_id 
            FROM posts p2 
            JOIN spots s2 ON p2.spot_id = s2.id 
            WHERE p2.user_id = $1
        )
        AND p.user_id != $1
        GROUP BY p.user_id`

	rows, err := r.db.QueryContext(ctx, query, userID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entities.ResonantUser
	for rows.Next() {
		var uid, count int
		if err := rows.Scan(&uid, &count); err != nil {
			return nil, err
		}
		idVO, _ := value_objects.NewID(uid)
		result = append(result, entities.ResonantUser{ID: idVO, MatchCount: count})
	}
	return result, nil
}

func (r *spotRepository) FindSpotsByMeshAndUsers(ctx context.Context, meshIDs []value_objects.MeshID, userIDs []value_objects.ID) ([]*entities.Spot, error) {
	query := `SELECT id, name, mesh_id, ST_X(location::geometry), ST_Y(location::geometry), registered_user_id 
              FROM spots WHERE mesh_id = ANY($1) AND registered_user_id = ANY($2)`

	mStrs := make([]string, len(meshIDs))
	for i, m := range meshIDs {
		mStrs[i] = m.String()
	}
	uInts := make([]int, len(userIDs))
	for i, u := range userIDs {
		uInts[i] = u.Value()
	}

	rows, err := r.db.QueryContext(ctx, query, pq.Array(mStrs), pq.Array(uInts))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var spots []*entities.Spot
	for rows.Next() {
		var sid, uid int
		var name, mID string
		var lng, lat float64
		if err := rows.Scan(&sid, &name, &mID, &lng, &lat, &uid); err != nil {
			return nil, err
		}
		s, _ := entities.NewSpot(sid, name, lat, lng, uid)
		spots = append(spots, s)
	}
	return spots, nil
}

func (r *spotRepository) GetDensityScoreByMesh(ctx context.Context, meshID value_objects.MeshID) (value_objects.DensityScore, error) {
	query := `SELECT count(*) FROM spots WHERE mesh_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, meshID.String()).Scan(&count)
	if err != nil {
		score, _ := value_objects.NewDensityScore(0)
		return score, err
	}
	score, _ := value_objects.NewDensityScore(count)
	return score, nil
}

func (r *spotRepository) FindPostsBySpot(ctx context.Context, spotID value_objects.ID) ([]*entities.Post, error) {
	query := `SELECT p.id, p.user_id, p.spot_id, p.username, p.image_url, p.caption, p.posted_at 
              FROM posts p 
              WHERE p.spot_id = $1`

	rows, err := r.db.QueryContext(ctx, query, spotID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*entities.Post
	for rows.Next() {
		var pid, uid, sid int
		var uname, img, capStr string
		var createdAt time.Time
		if err := rows.Scan(&pid, &uid, &sid, &uname, &img, &capStr, &createdAt); err != nil {
			return nil, err
		}
		p, _ := entities.NewPost(pid, uid, sid, uname, img, capStr, createdAt)
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *spotRepository) Update(spot *entities.Spot) error {
	query := `UPDATE spots SET name = $1, location = ST_SetSRID(ST_MakePoint($2, $3), 4326) WHERE id = $4`
	_, err := r.db.Exec(query, spot.Name.String(), spot.Longitude.Value(), spot.Latitude.Value(), spot.ID.Value())
	return err
}

func (r *spotRepository) Delete(id value_objects.ID) error {
	query := `DELETE FROM spots WHERE id = $1`
	_, err := r.db.Exec(query, id.Value())
	return err
}