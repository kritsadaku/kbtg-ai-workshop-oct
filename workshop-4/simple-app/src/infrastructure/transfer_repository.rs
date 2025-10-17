use async_trait::async_trait;
use sqlx::{SqlitePool, Row};
use chrono::Utc;
use uuid::Uuid;
use crate::domain::{
    Transfer, TransferRepository, CreateTransferRequest, TransferDb,
    PointLedger, PointLedgerRepository, EventType, PointLedgerDb,
};

#[derive(Clone)]
pub struct SqliteTransferRepository {
    pool: SqlitePool,
}

impl SqliteTransferRepository {
    pub fn new(pool: SqlitePool) -> Self {
        Self { pool }
    }

    pub async fn init_database(&self) -> Result<(), String> {
        // Create transfers table
        sqlx::query(
            r#"
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
            )
            "#,
        )
        .execute(&self.pool)
        .await
        .map_err(|e| format!("Failed to create transfers table: {}", e))?;

        // Create indexes
        sqlx::query("CREATE INDEX IF NOT EXISTS idx_transfers_from ON transfers(from_user_id)")
            .execute(&self.pool)
            .await
            .map_err(|e| format!("Failed to create index: {}", e))?;

        sqlx::query("CREATE INDEX IF NOT EXISTS idx_transfers_to ON transfers(to_user_id)")
            .execute(&self.pool)
            .await
            .map_err(|e| format!("Failed to create index: {}", e))?;

        sqlx::query("CREATE INDEX IF NOT EXISTS idx_transfers_created ON transfers(created_at)")
            .execute(&self.pool)
            .await
            .map_err(|e| format!("Failed to create index: {}", e))?;

        Ok(())
    }
}

#[async_trait]
impl TransferRepository for SqliteTransferRepository {
    async fn create_transfer(&self, transfer_request: CreateTransferRequest) -> Result<Transfer, String> {
        transfer_request.validate()?;
        
        let now = Utc::now();
        let idem_key = Uuid::new_v4().to_string();
        
        let result = sqlx::query(
            r#"
            INSERT INTO transfers (from_user_id, to_user_id, amount, status, note, idempotency_key, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
            "#,
        )
        .bind(transfer_request.from_user_id as i64)
        .bind(transfer_request.to_user_id as i64)
        .bind(transfer_request.amount as i64)
        .bind("pending")
        .bind(&transfer_request.note)
        .bind(&idem_key)
        .bind(now.to_rfc3339())
        .bind(now.to_rfc3339())
        .execute(&self.pool)
        .await
        .map_err(|e| format!("Failed to create transfer: {}", e))?;

        let transfer_id = result.last_insert_rowid() as u32;

        Ok(Transfer {
            idem_key,
            transfer_id: Some(transfer_id),
            from_user_id: transfer_request.from_user_id,
            to_user_id: transfer_request.to_user_id,
            amount: transfer_request.amount,
            status: crate::domain::TransferStatus::Pending,
            note: transfer_request.note,
            created_at: now,
            updated_at: now,
            completed_at: None,
            fail_reason: None,
        })
    }

    async fn get_transfer_by_idem_key(&self, idem_key: &str) -> Result<Option<Transfer>, String> {
        let row = sqlx::query(
            "SELECT id, from_user_id, to_user_id, amount, status, note, idempotency_key, created_at, updated_at, completed_at, fail_reason FROM transfers WHERE idempotency_key = ?"
        )
        .bind(idem_key)
        .fetch_optional(&self.pool)
        .await
        .map_err(|e| format!("Database error: {}", e))?;

        match row {
            Some(row) => {
                let transfer_db = TransferDb {
                    id: row.get::<i64, _>("id") as u32,
                    from_user_id: row.get::<i64, _>("from_user_id") as u32,
                    to_user_id: row.get::<i64, _>("to_user_id") as u32,
                    amount: row.get::<i64, _>("amount") as u32,
                    status: row.get("status"),
                    note: row.get("note"),
                    idempotency_key: row.get("idempotency_key"),
                    created_at: row.get("created_at"),
                    updated_at: row.get("updated_at"),
                    completed_at: row.get("completed_at"),
                    fail_reason: row.get("fail_reason"),
                };
                Ok(Some(transfer_db.to_domain()?))
            }
            None => Ok(None),
        }
    }

    async fn get_transfers_by_user_id(&self, user_id: u32, page: u32, page_size: u32) -> Result<(Vec<Transfer>, u32), String> {
        let limit = page_size as i64;
        let offset = ((page - 1) * page_size) as i64;

        // Get total count
        let total: i64 = sqlx::query_scalar(
            "SELECT COUNT(*) FROM transfers WHERE from_user_id = ? OR to_user_id = ?"
        )
        .bind(user_id as i64)
        .bind(user_id as i64)
        .fetch_one(&self.pool)
        .await
        .map_err(|e| format!("Database error: {}", e))?;

        // Get transfers
        let rows = sqlx::query(
            "SELECT id, from_user_id, to_user_id, amount, status, note, idempotency_key, created_at, updated_at, completed_at, fail_reason FROM transfers WHERE from_user_id = ? OR to_user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?"
        )
        .bind(user_id as i64)
        .bind(user_id as i64)
        .bind(limit)
        .bind(offset)
        .fetch_all(&self.pool)
        .await
        .map_err(|e| format!("Database error: {}", e))?;

        let mut transfers = Vec::new();
        for row in rows {
            let transfer_db = TransferDb {
                id: row.get::<i64, _>("id") as u32,
                from_user_id: row.get::<i64, _>("from_user_id") as u32,
                to_user_id: row.get::<i64, _>("to_user_id") as u32,
                amount: row.get::<i64, _>("amount") as u32,
                status: row.get("status"),
                note: row.get("note"),
                idempotency_key: row.get("idempotency_key"),
                created_at: row.get("created_at"),
                updated_at: row.get("updated_at"),
                completed_at: row.get("completed_at"),
                fail_reason: row.get("fail_reason"),
            };
            transfers.push(transfer_db.to_domain()?);
        }

        Ok((transfers, total as u32))
    }

    async fn update_transfer_status(&self, idem_key: &str, status: &str, completed_at: Option<String>, fail_reason: Option<String>) -> Result<(), String> {
        let now = Utc::now().to_rfc3339();
        
        sqlx::query(
            "UPDATE transfers SET status = ?, updated_at = ?, completed_at = ?, fail_reason = ? WHERE idempotency_key = ?"
        )
        .bind(status)
        .bind(now)
        .bind(completed_at)
        .bind(fail_reason)
        .bind(idem_key)
        .execute(&self.pool)
        .await
        .map_err(|e| format!("Failed to update transfer status: {}", e))?;

        Ok(())
    }
}

#[derive(Clone)]
pub struct SqlitePointLedgerRepository {
    pool: SqlitePool,
}

impl SqlitePointLedgerRepository {
    pub fn new(pool: SqlitePool) -> Self {
        Self { pool }
    }

    pub async fn init_database(&self) -> Result<(), String> {
        // Create point_ledger table
        sqlx::query(
            r#"
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
            )
            "#,
        )
        .execute(&self.pool)
        .await
        .map_err(|e| format!("Failed to create point_ledger table: {}", e))?;

        // Create indexes
        sqlx::query("CREATE INDEX IF NOT EXISTS idx_ledger_user ON point_ledger(user_id)")
            .execute(&self.pool)
            .await
            .map_err(|e| format!("Failed to create index: {}", e))?;

        sqlx::query("CREATE INDEX IF NOT EXISTS idx_ledger_transfer ON point_ledger(transfer_id)")
            .execute(&self.pool)
            .await
            .map_err(|e| format!("Failed to create index: {}", e))?;

        sqlx::query("CREATE INDEX IF NOT EXISTS idx_ledger_created ON point_ledger(created_at)")
            .execute(&self.pool)
            .await
            .map_err(|e| format!("Failed to create index: {}", e))?;

        Ok(())
    }
}

#[async_trait]
impl PointLedgerRepository for SqlitePointLedgerRepository {
    async fn create_ledger_entry(
        &self,
        user_id: u32,
        change: i32,
        balance_after: u32,
        event_type: EventType,
        transfer_id: Option<u32>,
        reference: Option<String>,
        metadata: Option<String>,
    ) -> Result<PointLedger, String> {
        let now = Utc::now();
        
        let result = sqlx::query(
            r#"
            INSERT INTO point_ledger (user_id, change, balance_after, event_type, transfer_id, reference, metadata, created_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)
            "#,
        )
        .bind(user_id as i64)
        .bind(change as i64)
        .bind(balance_after as i64)
        .bind(event_type.to_string())
        .bind(transfer_id.map(|id| id as i64))
        .bind(&reference)
        .bind(&metadata)
        .bind(now.to_rfc3339())
        .execute(&self.pool)
        .await
        .map_err(|e| format!("Failed to create ledger entry: {}", e))?;

        let id = result.last_insert_rowid() as u32;

        Ok(PointLedger {
            id,
            user_id,
            change,
            balance_after,
            event_type,
            transfer_id,
            reference,
            metadata,
            created_at: now,
        })
    }

    async fn get_ledger_by_user_id(&self, user_id: u32, limit: Option<i64>, offset: Option<i64>) -> Result<Vec<PointLedger>, String> {
        let limit = limit.unwrap_or(100);
        let offset = offset.unwrap_or(0);

        let rows = sqlx::query(
            "SELECT id, user_id, change, balance_after, event_type, transfer_id, reference, metadata, created_at FROM point_ledger WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?"
        )
        .bind(user_id as i64)
        .bind(limit)
        .bind(offset)
        .fetch_all(&self.pool)
        .await
        .map_err(|e| format!("Database error: {}", e))?;

        let mut ledger_entries = Vec::new();
        for row in rows {
            let ledger_db = PointLedgerDb {
                id: row.get::<i64, _>("id") as u32,
                user_id: row.get::<i64, _>("user_id") as u32,
                change: row.get::<i64, _>("change") as i32,
                balance_after: row.get::<i64, _>("balance_after") as u32,
                event_type: row.get("event_type"),
                transfer_id: row.get::<Option<i64>, _>("transfer_id").map(|id| id as u32),
                reference: row.get("reference"),
                metadata: row.get("metadata"),
                created_at: row.get("created_at"),
            };
            ledger_entries.push(ledger_db.to_domain()?);
        }

        Ok(ledger_entries)
    }

    async fn get_current_balance(&self, user_id: u32) -> Result<u32, String> {
        // Try to get the latest balance from point_ledger
        let balance: Option<i64> = sqlx::query_scalar(
            "SELECT balance_after FROM point_ledger WHERE user_id = ? ORDER BY created_at DESC LIMIT 1"
        )
        .bind(user_id as i64)
        .fetch_optional(&self.pool)
        .await
        .map_err(|e| format!("Database error: {}", e))?;

        // If no ledger entry exists, get from users table
        if let Some(balance) = balance {
            Ok(balance as u32)
        } else {
            let user_points: Option<i64> = sqlx::query_scalar(
                "SELECT points FROM users WHERE id = ?"
            )
            .bind(user_id as i64)
            .fetch_optional(&self.pool)
            .await
            .map_err(|e| format!("Database error: {}", e))?;

            Ok(user_points.unwrap_or(0) as u32)
        }
    }
}