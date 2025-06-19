package data

import (
	"database/sql"
	"secureGuard/internal/models"
)

type VulnerabilityModel struct {
	DB *sql.DB
}

func (m *VulnerabilityModel) Insert(vuln *models.Vulnerability) error {
	query := `
		INSERT INTO vulnerabilities (cve_id, title, description, severity, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	return m.DB.QueryRow(
		query,
		vuln.CveID,
		vuln.Title,
		vuln.Description,
		vuln.Severity,
		vuln.Status,
	).Scan(&vuln.Id, &vuln.CreatedAt, &vuln.UpdatedAt)
}

func (m *VulnerabilityModel) Update(vuln *models.Vulnerability) error {
	query := `
		UPDATE vulnerabilities
		SET cve_id = $1, title = $2, description = $3, severity = $4, status = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at`
	return m.DB.QueryRow(
		query,
		vuln.CveID,
		vuln.Title,
		vuln.Description,
		vuln.Severity,
		vuln.Status,
		vuln.Id,
	).Scan(&vuln.UpdatedAt)
}

func (m *VulnerabilityModel) Delete(id string) error {
	query := `DELETE FROM vulnerabilities WHERE id = $1`
	_, err := m.DB.Exec(query, id)
	return err
}

func (m *VulnerabilityModel) Get(id string) (*models.Vulnerability, error) {
	query := `
		SELECT id, cve_id, title, description, severity, status, created_at, updated_at
		FROM vulnerabilities
		WHERE id = $1`
	vuln := &models.Vulnerability{}
	err := m.DB.QueryRow(query, id).Scan(
		&vuln.Id,
		&vuln.CveID,
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
			&vuln.CveID,
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
	return vulns, rows.Err()
}

func (m *AssetModel) GetAssetsByVulnerability(vulnId string) ([]*models.Asset, error) {
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

func (m *VulnerabilityModel) AssociateAsset(vulnId, assetId string) error {
	query := `INSERT INTO asset_vulnerabilities (asset_id, vulnerability_id) VALUES ($1, $2)`
	_, err := m.DB.Exec(query, assetId, vulnId)
	return err
}
