package data

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"secureGuard/internal/models"
)

type IncidentModel struct {
	DB *sql.DB
}

func (m *IncidentModel) Insert(incident *models.Incident) error {
	if err := models.ValidateIncident(m.DB, incident); err != nil {
		return err
	}
	var assignedTo interface{}
	if incident.AssignedTo == "" {
		assignedTo = nil
	} else {
		assignedTo = incident.AssignedTo
	}
	_, err := m.DB.Exec(`INSERT INTO incidents (title,description,category,severity,status,reported_by,assigned_to)
							  VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		incident.Title, incident.Description, incident.Category,
		incident.Severity, incident.Status, incident.ReportedBy, assignedTo)
	return err
}

func (m *IncidentModel) Update(incident *models.Incident) error {
	current, err := m.Get(incident.Id)
	if err != nil {
		return err
	}
	if incident.Title == "" {
		incident.Title = current.Title
	}
	if incident.Description == "" {
		incident.Description = current.Description
	}
	if incident.Category == "" {
		incident.Category = current.Category
	}
	if incident.Severity == "" {
		incident.Severity = current.Severity
	}
	if incident.Status == "" {
		incident.Status = current.Status
	}
	if incident.ReportedBy == "" {
		incident.ReportedBy = current.ReportedBy
	}
	if incident.AssignedTo == "" {
		incident.AssignedTo = current.AssignedTo
	}
	if err := models.ValidateIncident(m.DB, incident); err != nil {
		return err
	}
	_, err = m.DB.Exec(
		`UPDATE incidents SET title = $1, description = $2, category = $3, severity = $4, status = $5, reported_by = $6, assigned_to = $7, updated_at = NOW() WHERE id = $8`,
		incident.Title,
		incident.Description,
		incident.Category,
		incident.Severity,
		incident.Status,
		incident.ReportedBy,
		incident.AssignedTo,
		incident.Id,
	)
	return err
}

func (m *IncidentModel) Delete(id string) error {
	result, err := m.DB.Exec(`DELETE FROM incidents WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("incident not found")
	}
	return nil
}

func (m *IncidentModel) Get(id string) (*models.Incident, error) {
	query := `
		SELECT id, title, description, category, severity, status, reported_by, assigned_to, created_at, updated_at
		FROM incidents
		WHERE id = $1`
	var AssignedTo sql.NullString
	incident := &models.Incident{}
	err := m.DB.QueryRow(query, id).Scan(
		&incident.Id,
		&incident.Title,
		&incident.Description,
		&incident.Category,
		&incident.Severity,
		&incident.Status,
		&incident.ReportedBy,
		&AssignedTo,
		&incident.CreatedAt,
		&incident.UpdatedAt,
	)
	if AssignedTo.Valid {
		incident.AssignedTo = AssignedTo.String
	} else {
		incident.AssignedTo = "null"
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("incident not found")
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
	var AssignedTo sql.NullString
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
			&AssignedTo,
			&incident.CreatedAt,
			&incident.UpdatedAt,
		)
		if AssignedTo.Valid {
			incident.AssignedTo = AssignedTo.String
		} else {
			incident.AssignedTo = "null"
		}
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, incident)
	}
	return incidents, nil
}

func (m *IncidentModel) GetAssetsByIncident(incidentId string) ([]*models.Asset, error) {
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
	query := `INSERT INTO incident_assets (incident_id, asset_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(query, incidentId, assetId)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return errors.New("association already exists")
		}
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return errors.New("incident or asset does not exist")
		}
		return err
	}
	return nil
}
