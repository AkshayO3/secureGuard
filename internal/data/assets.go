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
	if err := models.ValidateAsset(asset); err != nil {
		return err
	}
	_, err := m.DB.Exec(
		`INSERT INTO assets (name, type, ip_address, os, status) VALUES ($1, $2, $3, $4, $5)`,
		asset.Name,
		asset.Type,
		asset.Ip,
		asset.Os,
		asset.Status,
	)
	return err
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
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("asset not found")
	}
	return asset, err
}

func (m *AssetModel) Update(asset *models.Asset) error {
	current, err := m.Get(asset.Id)
	if err != nil {
		return err
	}
	if err := models.ValidateAsset(current); err != nil {
		return err
	}
	if asset.Name == "" {
		asset.Name = current.Name
	}
	if asset.Type == "" {
		asset.Type = current.Type
	}
	if asset.Ip == "" {
		asset.Ip = current.Ip
	}
	if asset.Os == "" {
		asset.Os = current.Os
	}
	if asset.Status == "" {
		asset.Status = current.Status
	}
	_, err = m.DB.Exec(
		`UPDATE assets SET name = $1, type = $2, ip_address = $3, os = $4, status = $5, updated_at = NOW() WHERE id = $6`,
		asset.Name, asset.Type, asset.Ip, asset.Os, asset.Status, asset.Id,
	)
	return err
}

func (m *AssetModel) Delete(id string) error {
	result, err := m.DB.Exec(`DELETE FROM assets WHERE id = $1`, id)
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
		FROM assets ORDER BY created_at DESC `
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
		SELECT v.id, v.cve_id,v.title, v.description, v.severity, v.status, v.created_at, v.updated_at
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
			&v.CveId,
			&v.Title,
			&v.Description,
			&v.Severity,
			&v.Status,
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
		FROM incident_assets ia
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
