# Database Schema Documentation

This document describes the database schema for the LBK Points Transfer API system, implemented using SQLite.

## Database Overview

The system consists of 3 main tables:
- **users**: User profile and point balance management
- **transfers**: Point transfer transactions between users
- **point_ledger**: Audit trail for all point transactions (append-only)

## Entity Relationship Diagram

```mermaid
erDiagram
    users {
        INTEGER id PK "Primary Key, Auto-increment"
        TEXT first_name "User's first name, NOT NULL"
        TEXT last_name "User's last name, NOT NULL"
        TEXT phone "User's phone number, NOT NULL"
        TEXT email "User's email address, UNIQUE, NOT NULL"
        TEXT member_since "Membership start date, NOT NULL"
        TEXT membership_level "Bronze/Silver/Gold/Platinum, NOT NULL"
        INTEGER points "Current point balance, DEFAULT 0"
        TEXT created_at "Record creation timestamp, NOT NULL"
        TEXT updated_at "Last update timestamp, NOT NULL"
    }

    transfers {
        INTEGER id PK "Primary Key, Auto-increment"
        INTEGER from_user_id FK "Sender user ID, NOT NULL"
        INTEGER to_user_id FK "Receiver user ID, NOT NULL"
        INTEGER amount "Transfer amount in points, CHECK > 0"
        TEXT status "Transfer status, CHECK IN enum, NOT NULL"
        TEXT note "Optional transfer note"
        TEXT idempotency_key "Unique key for duplicate prevention, UNIQUE"
        TEXT created_at "Transfer creation time, NOT NULL"
        TEXT updated_at "Last status update time, NOT NULL"
        TEXT completed_at "Transfer completion time, nullable"
        TEXT fail_reason "Failure reason if status=failed, nullable"
    }

    point_ledger {
        INTEGER id PK "Primary Key, Auto-increment"
        INTEGER user_id FK "User ID, NOT NULL"
        INTEGER change "Point change (+/-), NOT NULL"
        INTEGER balance_after "Balance after transaction, NOT NULL"
        TEXT event_type "Event type enum, CHECK IN values, NOT NULL"
        INTEGER transfer_id FK "Related transfer ID, nullable"
        TEXT reference "Transaction reference, nullable"
        TEXT metadata "JSON metadata, nullable"
        TEXT created_at "Transaction timestamp, NOT NULL"
    }

    %% Relationships
    users ||--o{ transfers : "sends (from_user_id)"
    users ||--o{ transfers : "receives (to_user_id)"
    users ||--o{ point_ledger : "has_transactions"
    transfers ||--o{ point_ledger : "generates_ledger_entries"
```

## Table Details

### users
**Purpose**: Store user profile information and current point balance.

**Key Features**:
- Auto-incrementing primary key
- Email uniqueness constraint
- Membership levels: Bronze, Silver, Gold, Platinum
- Point balance tracking
- Timestamp tracking for audit

**Sample Data**:
```sql
INSERT INTO users (first_name, last_name, phone, email, member_since, membership_level, points, created_at, updated_at) 
VALUES ('John', 'Doe', '+66812345678', 'john.doe@example.com', '2025-01-01T00:00:00Z', 'Gold', 1500, '2025-01-01T00:00:00Z', '2025-01-01T00:00:00Z');
```

### transfers
**Purpose**: Store point transfer transactions between users with status tracking.

**Key Features**:
- Bidirectional foreign keys to users table
- Status enum validation: 'pending', 'processing', 'completed', 'failed', 'cancelled', 'reversed'
- Idempotency key for duplicate prevention
- Amount validation (must be > 0)
- Complete audit trail with timestamps

**Indexes**:
- `idx_transfers_from` on from_user_id
- `idx_transfers_to` on to_user_id  
- `idx_transfers_created` on created_at

### point_ledger
**Purpose**: Append-only audit log for all point transactions.

**Key Features**:
- Event type validation: 'transfer_out', 'transfer_in', 'adjust', 'earn', 'redeem'
- Balance tracking after each transaction
- Optional link to transfer record
- JSON metadata support for additional context
- Complete transaction history preservation

**Indexes**:
- `idx_ledger_user` on user_id
- `idx_ledger_transfer` on transfer_id
- `idx_ledger_created` on created_at

## Business Rules

1. **Transfer Atomicity**: Each transfer creates exactly 2 ledger entries (one debit, one credit)
2. **Balance Consistency**: Point ledger maintains running balance for each user
3. **Idempotency**: Duplicate transfers prevented by unique idempotency_key
4. **Audit Trail**: All point movements recorded in point_ledger (append-only)
5. **Status Flow**: Transfers progress: pending → (processing) → completed/failed

## Relationships

- **One-to-Many**: users → transfers (via from_user_id)
- **One-to-Many**: users → transfers (via to_user_id)  
- **One-to-Many**: users → point_ledger (transaction history)
- **One-to-Many**: transfers → point_ledger (each transfer generates 2 entries)

## Data Integrity

### Foreign Key Constraints
- `transfers.from_user_id` → `users.id`
- `transfers.to_user_id` → `users.id`
- `point_ledger.user_id` → `users.id`
- `point_ledger.transfer_id` → `transfers.id`

### Check Constraints
- `transfers.amount > 0`
- `transfers.status IN ('pending','processing','completed','failed','cancelled','reversed')`
- `point_ledger.event_type IN ('transfer_out','transfer_in','adjust','earn','redeem')`

### Unique Constraints
- `users.email` (unique across system)
- `transfers.idempotency_key` (prevents duplicate transfers)

## Performance Considerations

- Indexes on frequently queried columns (user_id, created_at, transfer relationships)
- Append-only ledger design for optimal insert performance
- Separate balance tracking in ledger prevents need for complex aggregations
- Idempotency key lookup optimization for duplicate detection