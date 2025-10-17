use async_trait::async_trait;
use super::user::{User, CreateUserRequest, UpdateUserRequest};
use super::transfer::{Transfer, CreateTransferRequest, TransferDb};
use super::point_ledger::{PointLedger, EventType, PointLedgerDb};

#[async_trait]
pub trait UserRepository {
    async fn get_user_by_id(&self, id: u32) -> Result<Option<User>, String>;
    async fn get_user_by_email(&self, email: &str) -> Result<Option<User>, String>;
    async fn create_user(&self, user_request: CreateUserRequest) -> Result<User, String>;
    async fn update_user(&self, id: u32, update_request: UpdateUserRequest) -> Result<User, String>;
    async fn delete_user(&self, id: u32) -> Result<bool, String>;
    async fn list_users(&self, limit: Option<i64>, offset: Option<i64>) -> Result<Vec<User>, String>;
}

#[async_trait]
pub trait TransferRepository {
    async fn create_transfer(&self, transfer_request: CreateTransferRequest) -> Result<Transfer, String>;
    async fn get_transfer_by_idem_key(&self, idem_key: &str) -> Result<Option<Transfer>, String>;
    async fn get_transfers_by_user_id(&self, user_id: u32, page: u32, page_size: u32) -> Result<(Vec<Transfer>, u32), String>;
    async fn update_transfer_status(&self, idem_key: &str, status: &str, completed_at: Option<String>, fail_reason: Option<String>) -> Result<(), String>;
}

#[async_trait]
pub trait PointLedgerRepository {
    async fn create_ledger_entry(&self, user_id: u32, change: i32, balance_after: u32, event_type: EventType, transfer_id: Option<u32>, reference: Option<String>, metadata: Option<String>) -> Result<PointLedger, String>;
    async fn get_ledger_by_user_id(&self, user_id: u32, limit: Option<i64>, offset: Option<i64>) -> Result<Vec<PointLedger>, String>;
    async fn get_current_balance(&self, user_id: u32) -> Result<u32, String>;
}