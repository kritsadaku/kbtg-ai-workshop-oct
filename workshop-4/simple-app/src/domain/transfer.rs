use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use utoipa::ToSchema;

#[derive(Debug, Clone, Serialize, Deserialize, ToSchema)]
#[serde(rename_all = "lowercase")]
pub enum TransferStatus {
    Pending,
    Processing,
    Completed,
    Failed,
    Cancelled,
    Reversed,
}

impl std::fmt::Display for TransferStatus {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            TransferStatus::Pending => write!(f, "pending"),
            TransferStatus::Processing => write!(f, "processing"),
            TransferStatus::Completed => write!(f, "completed"),
            TransferStatus::Failed => write!(f, "failed"),
            TransferStatus::Cancelled => write!(f, "cancelled"),
            TransferStatus::Reversed => write!(f, "reversed"),
        }
    }
}

impl std::str::FromStr for TransferStatus {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s.to_lowercase().as_str() {
            "pending" => Ok(TransferStatus::Pending),
            "processing" => Ok(TransferStatus::Processing),
            "completed" => Ok(TransferStatus::Completed),
            "failed" => Ok(TransferStatus::Failed),
            "cancelled" => Ok(TransferStatus::Cancelled),
            "reversed" => Ok(TransferStatus::Reversed),
            _ => Err(format!("Invalid transfer status: {}", s)),
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, ToSchema)]
pub struct Transfer {
    #[serde(rename = "idemKey")]
    pub idem_key: String,
    #[serde(rename = "transferId")]
    pub transfer_id: Option<u32>,
    #[serde(rename = "fromUserId")]
    pub from_user_id: u32,
    #[serde(rename = "toUserId")]
    pub to_user_id: u32,
    pub amount: u32,
    pub status: TransferStatus,
    pub note: Option<String>,
    #[serde(rename = "createdAt")]
    #[schema(value_type = String, format = "date-time")]
    pub created_at: DateTime<Utc>,
    #[serde(rename = "updatedAt")]
    #[schema(value_type = String, format = "date-time")]
    pub updated_at: DateTime<Utc>,
    #[serde(rename = "completedAt")]
    #[schema(value_type = Option<String>, format = "date-time")]
    pub completed_at: Option<DateTime<Utc>>,
    #[serde(rename = "failReason")]
    pub fail_reason: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, ToSchema)]
pub struct CreateTransferRequest {
    #[serde(rename = "fromUserId")]
    pub from_user_id: u32,
    #[serde(rename = "toUserId")]
    pub to_user_id: u32,
    pub amount: u32,
    pub note: Option<String>,
}

impl CreateTransferRequest {
    pub fn validate(&self) -> Result<(), String> {
        if self.amount == 0 {
            return Err("Amount must be greater than 0".to_string());
        }
        
        if self.from_user_id == self.to_user_id {
            return Err("Cannot transfer to the same user".to_string());
        }
        
        if let Some(note) = &self.note {
            if note.len() > 512 {
                return Err("Note cannot exceed 512 characters".to_string());
            }
        }
        
        Ok(())
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, ToSchema)]
pub struct TransferCreateResponse {
    pub transfer: Transfer,
}

#[derive(Debug, Clone, Serialize, Deserialize, ToSchema)]
pub struct TransferGetResponse {
    pub transfer: Transfer,
}

#[derive(Debug, Clone, Serialize, Deserialize, ToSchema)]
pub struct TransferListResponse {
    pub data: Vec<Transfer>,
    pub page: u32,
    #[serde(rename = "pageSize")]
    pub page_size: u32,
    pub total: u32,
}

// Database model for internal use
#[derive(Debug, Clone)]
pub struct TransferDb {
    pub id: u32,
    pub from_user_id: u32,
    pub to_user_id: u32,
    pub amount: u32,
    pub status: String,
    pub note: Option<String>,
    pub idempotency_key: String,
    pub created_at: String,
    pub updated_at: String,
    pub completed_at: Option<String>,
    pub fail_reason: Option<String>,
}

impl TransferDb {
    pub fn to_domain(self) -> Result<Transfer, String> {
        let status = self.status.parse::<TransferStatus>()
            .map_err(|e| format!("Invalid status: {}", e))?;
        
        let created_at = DateTime::parse_from_rfc3339(&self.created_at)
            .map_err(|e| format!("Invalid created_at date: {}", e))?
            .with_timezone(&Utc);
        
        let updated_at = DateTime::parse_from_rfc3339(&self.updated_at)
            .map_err(|e| format!("Invalid updated_at date: {}", e))?
            .with_timezone(&Utc);
        
        let completed_at = if let Some(completed_str) = self.completed_at {
            Some(DateTime::parse_from_rfc3339(&completed_str)
                .map_err(|e| format!("Invalid completed_at date: {}", e))?
                .with_timezone(&Utc))
        } else {
            None
        };
        
        Ok(Transfer {
            idem_key: self.idempotency_key,
            transfer_id: Some(self.id),
            from_user_id: self.from_user_id,
            to_user_id: self.to_user_id,
            amount: self.amount,
            status,
            note: self.note,
            created_at,
            updated_at,
            completed_at,
            fail_reason: self.fail_reason,
        })
    }
}