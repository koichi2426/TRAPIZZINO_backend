package postgres

import (
	"database/sql"
	"app/domain/entities"
	"app/domain/value_objects"
)

// SpotRepositoryはPostgreSQL+PostGISを用いたSpotRepositoryの実装です。
type SpotRepository struct {
	db *sql.DB
}

func NewSpotRepository(db *sql.DB) *SpotRepository {
	return &SpotRepository{db: db}
}

func (r *SpotRepository) Create(spot *entities.Spot) (*entities.Spot, error) {
	query := `INSERT INTO spots (name, mesh_id, location) VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326))
	ON CONFLICT (mesh_id) DO UPDATE SET name = EXCLUDED.name, location = EXCLUDED.location RETURNING id`
	var id int
	err := r.db.QueryRow(query, spot.Name, spot.MeshID.String(), spot.Longitude.Value(), spot.Latitude.Value()).Scan(&id)
	if err != nil {
		return nil, err
	}
	spot.ID, _ = value_objects.NewID(id)
	return spot, nil
}

func (r *SpotRepository) FindByID(id value_objects.ID) (*entities.Spot, error) {
	query := `SELECT id, name, mesh_id, ST_X(location::geometry), ST_Y(location::geometry) FROM spots WHERE id = $1`
	row := r.db.QueryRow(query, id.Value())
	var sid int
	var name, meshID string
	var lng, lat float64
	if err := row.Scan(&sid, &name, &meshID, &lng, &lat); err != nil {
		return nil, err
	}
	spotID, _ := value_objects.NewID(sid)
	mesh, _ := value_objects.NewMeshID(lat, lng)
	latitude, _ := value_objects.NewLatitude(lat)
	longitude, _ := value_objects.NewLongitude(lng)
	spotName, _ := value_objects.NewSpotName(name)
	return &entities.Spot{
		ID:        spotID,
		Name:      spotName,
		MeshID:    mesh,
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}

func (r *SpotRepository) FindByMeshID(meshID value_objects.MeshID) ([]*entities.Spot, error) {
	query := `SELECT id, name, mesh_id, ST_X(location::geometry), ST_Y(location::geometry) FROM spots WHERE mesh_id = $1`
	rows, err := r.db.Query(query, meshID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var spots []*entities.Spot
	for rows.Next() {
		var sid int
		var name, mesh string
		var lng, lat float64
		if err := rows.Scan(&sid, &name, &mesh, &lng, &lat); err != nil {
			return nil, err
		}
		spotID, _ := value_objects.NewID(sid)
		meshVO, _ := value_objects.NewMeshID(lat, lng)
		latitude, _ := value_objects.NewLatitude(lat)
		longitude, _ := value_objects.NewLongitude(lng)
		spotName, _ := value_objects.NewSpotName(name)
		spots = append(spots, &entities.Spot{
			ID:        spotID,
			Name:      spotName,
			MeshID:    meshVO,
			Latitude:  latitude,
			Longitude: longitude,
		})
	}
	return spots, nil
}

func (r *SpotRepository) Update(spot *entities.Spot) error {
	query := `UPDATE spots SET name = $1, location = ST_SetSRID(ST_MakePoint($2, $3), 4326) WHERE id = $4`
	_, err := r.db.Exec(query, spot.Name, spot.Longitude.Value(), spot.Latitude.Value(), spot.ID.Value())
	return err
}

func (r *SpotRepository) Delete(id value_objects.ID) error {
	query := `DELETE FROM spots WHERE id = $1`
	_, err := r.db.Exec(query, id.Value())
	return err
}
