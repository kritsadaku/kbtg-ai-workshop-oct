package repository

import (
	"database/sql"
	"temp-backend/internal/domain"
	"time"
)

// pointLedgerRepository implements domain.PointLedgerRepository
type pointLedgerRepository struct {
	db *sql.DB
}

// NewPointLedgerRepository creates a new point ledger repository
func NewPointLedgerRepository(db *sql.DB) domain.PointLedgerRepository {
	return &pointLedgerRepository{
		db: db,
	}
}

// Create creates a new point ledger entry
func (r *pointLedgerRepository) Create(entry *domain.PointLedger) error {
	now := time.Now().Format("2006-01-02 15:04:05")

	_, err := r.db.Exec(`
		INSERT INTO point_ledger (user_id, change, balance_after, event_type, transfer_id, reference, metadata, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.UserID, entry.Change, entry.BalanceAfter, entry.EventType,
		entry.TransferID, entry.Reference, entry.Metadata, now)

	return err
}

// GetByUserID retrieves point ledger entries by user ID
func (r *pointLedgerRepository) GetByUserID(userID int, limit int) ([]*domain.PointLedger, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, change, balance_after, event_type, transfer_id, reference, metadata, created_at
		FROM point_ledger 
		WHERE user_id = ?
		ORDER BY created_at DESC
		LIMIT ?`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*domain.PointLedger
	for rows.Next() {
		entry := &domain.PointLedger{}
		var transferID sql.NullInt64
		var reference, metadata sql.NullString

		err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.Change, &entry.BalanceAfter,
			&entry.EventType, &transferID, &reference, &metadata, &entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if transferID.Valid {
			transferIDInt := int(transferID.Int64)
			entry.TransferID = &transferIDInt
		}
		if reference.Valid {
			entry.Reference = &reference.String
		}
		if metadata.Valid {
			entry.Metadata = &metadata.String
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
