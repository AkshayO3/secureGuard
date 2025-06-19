package data

import "database/sql"

type IncidentModel struct {
	DB *sql.DB
}

func (m *IncidentModel) Insert(incident *models.Incident) error {
	query := `
		INSERT INTO incidents (title, description, category, severity, status, reported_by, assigned_to)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`
	return m.DB.QueryRow(
		query,
		incident.Title,
		incident.Description,
		incident.Category,
		incident.Severity,
		incident.Status,
		incident.ReportedBy,
		incident.AssignedTo,
	).Scan(&incident.Id, &incident.CreatedAt, &incident.UpdatedAt)
}

func (m *IncidentModel) Update(incident *models.Incident) error {
	query := `
		UPDATE incidents
		SET title = $1, description = $2, category = $3, severity = $4, status = $5, reported_by = $6, assigned_to = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at`
	return m.DB.QueryRow(
		query,
		incident.Title,
		incident.Description,
		incident.Category,
		incident.Severity,
		incident.Status,
		incident.ReportedBy,
		incident.AssignedTo,
		incident.Id,
	).Scan(&incident.UpdatedAt)
}

func (m *IncidentModel) Delete(id string) error {
	query := `DELETE FROM incidents WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}

func (m *IncidentModel) Get(id string) (*models.Incident, error) {
	query := `
		SELECT id, title, description, category, severity, status, reported_by, assigned_to, created_at, updated_at
		FROM incidents
		WHERE id = $1`
	incident := &models.Incident{}
	err := m.DB.QueryRow(query, id).Scan(
		&incident.Id,
		&incident.Title,
		&incident.Description,
		&incident.Category,
		&incident.Severity,
		&incident.Status,
		&incident.ReportedBy,
		&incident.AssignedTo,
		&incident.CreatedAt,
		&incident.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return incident, nil
}

func (m *IncidentModel) List() ([]*models.Incident, error) {
	query := `
		SELECT id, title, description, category, severity, status, reported_by, assigned_to, created_at, updated_at
		FROM incidents
		ORDER BY created_at DESC`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []*models.Incident
	for rows.Next() {
		incident := &models.Incident{}
		err := rows.Scan(
			&incident.Id,
			&incident.Title,
			&incident.Description,
			&incident.Category,
			&incident.Severity,
			&incident.Status,
			&incident.ReportedBy,
			&incident.AssignedTo,
			&incident.CreatedAt,
			&incident.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, incident)
	}
	return incidents, rows.Err()
}

func (m *AssetModel) GetAssetsByIncident(incidentId string) ([]*models.Asset, error) {
	query := `
		SELECT a.id, a.name, a.type, a.ip_address, a.os, a.status, a.created_at, a.updated_at
		FROM incidentassets ia
		JOIN assets a ON ia.asset_id = a.id
		WHERE ia.incident_id = $1;`
	rows, err := m.DB.Query(query, incidentId)
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

func (m *IncidentModel) AssociateAsset(incidentId, assetId string) error {
	query := `INSERT INTO incidentassets (incident_id, asset_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(query, incidentId, assetId)
	return err
}
