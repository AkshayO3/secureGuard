package data

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"secureGuard/internal/models"
)

type VulnerabilityModel struct {
	DB *sql.DB
}

func (m *VulnerabilityModel) Insert(vuln *models.Vulnerability) error {
	if err := models.ValidateVulnerability(vuln); err != nil {
		return err
	}
	_, err := m.DB.Exec(
		`INSERT INTO vulnerabilities (cve_id, title, description, severity, status) VALUES ($1, $2, $3, $4, $5)`,
		vuln.CveId,
		vuln.Title,
		vuln.Description,
		vuln.Severity,
		vuln.Status,
	)
	return err
}

func (m *VulnerabilityModel) Update(vuln *models.Vulnerability) error {
	current, err := m.Get(vuln.Id)
	if err != nil {
		return err
	}
	if err := models.ValidateVulnerability(current); err != nil {
		return err
	}
	if vuln.CveId == "" {
		vuln.CveId = current.CveId
	}
	if vuln.Title == "" {
		vuln.Title = current.Title
	}
	if vuln.Description == "" {
		vuln.Description = current.Description
	}
	if vuln.Severity == "" {
		vuln.Severity = current.Severity
	}
	if vuln.Status == "" {
		vuln.Status = current.Status
	}
	_, err = m.DB.Exec(
		`UPDATE vulnerabilities SET cve_id = $1, title = $2, description = $3, severity = $4, status = $5, updated_at = NOW() WHERE id = $6`,
		vuln.CveId,
		vuln.Title,
		vuln.Description,
		vuln.Severity,
		vuln.Status,
		vuln.Id,
	)
	return err
}

func (m *VulnerabilityModel) Delete(id string) error {
	result, err := m.DB.Exec(`DELETE FROM vulnerabilities WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("vulnerability not found")
	}
	return nil
}

func (m *VulnerabilityModel) Get(id string) (*models.Vulnerability, error) {
	query := `
		SELECT id, cve_id, title, description, severity, status, created_at, updated_at
		FROM vulnerabilities
		WHERE id = $1`
	vuln := &models.Vulnerability{}
	err := m.DB.QueryRow(query, id).Scan(
		&vuln.Id,
		&vuln.CveId,
		&vuln.Title,
		&vuln.Description,
		&vuln.Severity,
		&vuln.Status,
		&vuln.CreatedAt,
		&vuln.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return vuln, nil
}

func (m *VulnerabilityModel) List() ([]*models.Vulnerability, error) {
	query := `
		SELECT id, cve_id, title, description, severity, status, created_at, updated_at
		FROM vulnerabilities
		ORDER BY created_at DESC`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vulns []*models.Vulnerability
	for rows.Next() {
		vuln := &models.Vulnerability{}
		err := rows.Scan(
			&vuln.Id,
			&vuln.CveId,
			&vuln.Title,
			&vuln.Description,
			&vuln.Severity,
			&vuln.Status,
			&vuln.CreatedAt,
			&vuln.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vulns = append(vulns, vuln)
	}
	return vulns, nil
}

func (m *VulnerabilityModel) GetAssetsByVulnerability(vulnId string) ([]*models.Asset, error) {
	query := `
		SELECT a.id, a.name, a.type, a.ip_address, a.os, a.status, a.created_at, a.updated_at
		FROM asset_vulnerabilities av
		JOIN assets a ON av.asset_id = a.id
		WHERE av.vulnerability_id = $1;`
	rows, err := m.DB.Query(query, vulnId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []*models.Asset
	for rows.Next() {
		a := &models.Asset{}
		err := rows.Scan(
			&a.Id,
			&a.Name,
			&a.Type,
			&a.Ip,
			&a.Os,
			&a.Status,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, nil
}

func (m *VulnerabilityModel) AssociateAsset(vulnId, assetId, status string) error {
	query := `INSERT INTO asset_vulnerabilities (asset_id, vulnerability_id, status) VALUES ($1, $2, $3)`
	_, err := m.DB.Exec(query, assetId, vulnId, status)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errors.New("asset-vulnerability association already exists")
		}
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return errors.New("asset or vulnerability does not exist")
		}
		return err
	}
	return nil
}
