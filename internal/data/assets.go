package data

import (
	"database/sql"
	"errors"
	"secureGuard/internal/models"
)

type AssetModel struct {
	DB *sql.DB
}

func (m *AssetModel) Insert(asset *models.Asset) error {
	query := `
		INSERT INTO assets (name, type, ip_address, os, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	return m.DB.QueryRow(
		query,
		asset.Name,
		asset.Type,
		asset.Ip,
		asset.Os,
		asset.Status,
	).Scan(&asset.Id, &asset.CreatedAt, &asset.UpdatedAt)
}

func (m *AssetModel) Get(id string) (*models.Asset, error) {
	query := `
		SELECT id, name, type, ip_address, os, status, created_at, updated_at
		FROM assets WHERE id = $1`
	asset := &models.Asset{}
	err := m.DB.QueryRow(query, id).Scan(
		&asset.Id,
		&asset.Name,
		&asset.Type,
		&asset.Ip,
		&asset.Os,
		&asset.Status,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("asset not found")
	}
	return asset, err
}

func (m *AssetModel) Update(asset *models.Asset) error {
	query := `
		UPDATE assets
		SET name = $1, type = $2, ip_address = $3, os = $4, status = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at`
	return m.DB.QueryRow(
		query,
		asset.Name,
		asset.Type,
		asset.Ip,
		asset.Os,
		asset.Status,
		asset.Id,
	).Scan(&asset.UpdatedAt)
}

func (m *AssetModel) Delete(id string) error {
	query := `DELETE FROM assets WHERE id = $1`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("asset not found")
	}
	return nil
}

func (m *AssetModel) List() ([]*models.Asset, error) {
	query := `
		SELECT id, name, type, ip_address, os, status, created_at, updated_at
		FROM assets`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*models.Asset
	for rows.Next() {
		asset := &models.Asset{}
		err := rows.Scan(
			&asset.Id,
			&asset.Name,
			&asset.Type,
			&asset.Ip,
			&asset.Os,
			&asset.Status,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return assets, nil
}

func (m *AssetModel) GetVulnerabilities(id string) ([]*models.Vulnerability, error) {
	query := `
		SELECT v.id, v.name, v.description, v.severity, v.created_at, v.updated_at
		FROM asset_vulnerabilities av
		JOIN vulnerabilities v ON av.vulnerability_id = v.id
		WHERE av.asset_id = $1;`
	rows, err := m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vulns []*models.Vulnerability
	for rows.Next() {
		v := &models.Vulnerability{}
		err := rows.Scan(
			&v.Id,
			&v.Name,
			&v.Description,
			&v.Severity,
			&v.CreatedAt,
			&v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vulns = append(vulns, v)
	}
	return vulns, nil
}

func (m *AssetModel) GetIncidents(id string) ([]*models.Incident, error) {
	query := `
		SELECT i.id, i.title, i.description, i.severity, i.status, i.created_at, i.updated_at
		FROM incidentassets ia
		JOIN incidents i ON ia.incident_id = i.id
		WHERE ia.asset_id = $1;`
	rows, err := m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []*models.Incident
	for rows.Next() {
		inc := &models.Incident{}
		err := rows.Scan(
			&inc.Id,
			&inc.Title,
			&inc.Description,
			&inc.Severity,
			&inc.Status,
			&inc.CreatedAt,
			&inc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, inc)
	}
	return incidents, nil
}
