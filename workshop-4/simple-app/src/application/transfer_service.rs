use std::sync::Arc;
use chrono::Utc;
use crate::domain::{
    Transfer, TransferRepository, CreateTransferRequest, TransferCreateResponse, TransferGetResponse, TransferListResponse,
    PointLedgerRepository, EventType, UserRepository,
};

#[derive(Clone)]
pub struct TransferService {
    transfer_repository: Arc<dyn TransferRepository + Send + Sync>,
    point_ledger_repository: Arc<dyn PointLedgerRepository + Send + Sync>,
    user_repository: Arc<dyn UserRepository + Send + Sync>,
}

impl TransferService {
    pub fn new(
        transfer_repository: Arc<dyn TransferRepository + Send + Sync>,
        point_ledger_repository: Arc<dyn PointLedgerRepository + Send + Sync>,
        user_repository: Arc<dyn UserRepository + Send + Sync>,
    ) -> Self {
        Self {
            transfer_repository,
            point_ledger_repository,
            user_repository,
        }
    }

    pub async fn create_transfer(&self, request: CreateTransferRequest) -> Result<TransferCreateResponse, String> {
        // Validate request
        request.validate()?;

        // Check if users exist
        let from_user = self.user_repository.get_user_by_id(request.from_user_id).await?
            .ok_or("From user not found".to_string())?;
        
        let _to_user = self.user_repository.get_user_by_id(request.to_user_id).await?
            .ok_or("To user not found".to_string())?;

        // Check if sender has enough points
        let current_balance = self.point_ledger_repository.get_current_balance(request.from_user_id).await?;
        if current_balance < request.amount {
            return Err("Insufficient points".to_string());
        }

        // Create transfer (initially pending)
        let mut transfer = self.transfer_repository.create_transfer(request.clone()).await?;

        // Process the transfer immediately (in a real system, this might be async)
        match self.process_transfer(&transfer).await {
            Ok(_) => {
                // Update transfer status to completed
                let completed_at = Some(Utc::now().to_rfc3339());
                self.transfer_repository.update_transfer_status(
                    &transfer.idem_key,
                    "completed",
                    completed_at.clone(),
                    None,
                ).await?;

                // Update transfer object
                transfer.status = crate::domain::TransferStatus::Completed;
                transfer.completed_at = completed_at.map(|s| s.parse().unwrap());
                transfer.updated_at = Utc::now();
            }
            Err(e) => {
                // Update transfer status to failed
                self.transfer_repository.update_transfer_status(
                    &transfer.idem_key,
                    "failed",
                    None,
                    Some(e.clone()),
                ).await?;

                // Update transfer object
                transfer.status = crate::domain::TransferStatus::Failed;
                transfer.fail_reason = Some(e);
                transfer.updated_at = Utc::now();
            }
        }

        Ok(TransferCreateResponse { transfer })
    }

    pub async fn get_transfer(&self, idem_key: &str) -> Result<TransferGetResponse, String> {
        let transfer = self.transfer_repository.get_transfer_by_idem_key(idem_key).await?
            .ok_or("Transfer not found".to_string())?;

        Ok(TransferGetResponse { transfer })
    }

    pub async fn list_transfers(&self, user_id: u32, page: u32, page_size: u32) -> Result<TransferListResponse, String> {
        // Validate parameters
        if page == 0 {
            return Err("Page must be greater than 0".to_string());
        }
        if page_size == 0 || page_size > 200 {
            return Err("Page size must be between 1 and 200".to_string());
        }

        // Check if user exists
        let _user = self.user_repository.get_user_by_id(user_id).await?
            .ok_or("User not found".to_string())?;

        let (transfers, total) = self.transfer_repository.get_transfers_by_user_id(user_id, page, page_size).await?;

        Ok(TransferListResponse {
            data: transfers,
            page,
            page_size,
            total,
        })
    }

    async fn process_transfer(&self, transfer: &Transfer) -> Result<(), String> {
        // Get current balances
        let from_balance = self.point_ledger_repository.get_current_balance(transfer.from_user_id).await?;
        let to_balance = self.point_ledger_repository.get_current_balance(transfer.to_user_id).await?;

        // Double-check sender has enough points
        if from_balance < transfer.amount {
            return Err("Insufficient points".to_string());
        }

        // Create ledger entries (transfer_out for sender)
        let new_from_balance = from_balance - transfer.amount;
        self.point_ledger_repository.create_ledger_entry(
            transfer.from_user_id,
            -(transfer.amount as i32),
            new_from_balance,
            EventType::TransferOut,
            transfer.transfer_id,
            Some(format!("Transfer to user {}", transfer.to_user_id)),
            Some(serde_json::json!({
                "transfer_id": transfer.transfer_id,
                "idem_key": transfer.idem_key,
                "note": transfer.note
            }).to_string()),
        ).await?;

        // Create ledger entry (transfer_in for receiver)
        let new_to_balance = to_balance + transfer.amount;
        self.point_ledger_repository.create_ledger_entry(
            transfer.to_user_id,
            transfer.amount as i32,
            new_to_balance,
            EventType::TransferIn,
            transfer.transfer_id,
            Some(format!("Transfer from user {}", transfer.from_user_id)),
            Some(serde_json::json!({
                "transfer_id": transfer.transfer_id,
                "idem_key": transfer.idem_key,
                "note": transfer.note
            }).to_string()),
        ).await?;

        Ok(())
    }
}