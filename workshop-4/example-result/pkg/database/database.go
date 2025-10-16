package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Database wraps the sql.DB connection
type Database struct {
	DB *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	database := &Database{DB: db}
	if err := database.migrate(); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return database, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}

// migrate creates the necessary tables
func (d *Database) migrate() error {
	// Users table
	createUsersSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		phone TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		member_since TEXT NOT NULL,
		membership_level TEXT NOT NULL DEFAULT 'Bronze',
		points INTEGER NOT NULL DEFAULT 0,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);`

	if _, err := d.DB.Exec(createUsersSQL); err != nil {
		return err
	}

	// Transfers table
	createTransfersSQL := `
	CREATE TABLE IF NOT EXISTS transfers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_user_id INTEGER NOT NULL,
		to_user_id INTEGER NOT NULL,
		amount INTEGER NOT NULL CHECK (amount > 0),
		status TEXT NOT NULL CHECK (status IN ('pending','processing','completed','failed','cancelled','reversed')),
		note TEXT,
		idempotency_key TEXT NOT NULL UNIQUE,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		completed_at TEXT,
		fail_reason TEXT,
		FOREIGN KEY (from_user_id) REFERENCES users(id),
		FOREIGN KEY (to_user_id) REFERENCES users(id)
	);`

	if _, err := d.DB.Exec(createTransfersSQL); err != nil {
		return err
	}

	// Point ledger table
	createPointLedgerSQL := `
	CREATE TABLE IF NOT EXISTS point_ledger (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		change INTEGER NOT NULL,
		balance_after INTEGER NOT NULL,
		event_type TEXT NOT NULL CHECK (event_type IN ('transfer_out','transfer_in','adjust','earn','redeem')),
		transfer_id INTEGER,
		reference TEXT,
		metadata TEXT,
		created_at TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (transfer_id) REFERENCES transfers(id)
	);`

	if _, err := d.DB.Exec(createPointLedgerSQL); err != nil {
		return err
	}

	// Create indexes for transfers
	createTransferIndexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_transfers_from ON transfers(from_user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_transfers_to ON transfers(to_user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_transfers_created ON transfers(created_at);`,
		`CREATE INDEX IF NOT EXISTS idx_transfers_idem_key ON transfers(idempotency_key);`,
	}

	for _, indexSQL := range createTransferIndexes {
		if _, err := d.DB.Exec(indexSQL); err != nil {
			return err
		}
	}

	// Create indexes for point ledger
	createLedgerIndexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_ledger_user ON point_ledger(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_ledger_transfer ON point_ledger(transfer_id);`,
		`CREATE INDEX IF NOT EXISTS idx_ledger_created ON point_ledger(created_at);`,
	}

	for _, indexSQL := range createLedgerIndexes {
		if _, err := d.DB.Exec(indexSQL); err != nil {
			return err
		}
	}

	return nil
}
