use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use utoipa::ToSchema;

#[derive(Debug, Clone, Serialize, Deserialize, ToSchema)]
#[serde(rename_all = "snake_case")]
pub enum EventType {
    TransferOut,
    TransferIn,
    Adjust,
    Earn,
    Redeem,
}

impl std::fmt::Display for EventType {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            EventType::TransferOut => write!(f, "transfer_out"),
            EventType::TransferIn => write!(f, "transfer_in"),
            EventType::Adjust => write!(f, "adjust"),
            EventType::Earn => write!(f, "earn"),
            EventType::Redeem => write!(f, "redeem"),
        }
    }
}

impl std::str::FromStr for EventType {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s {
            "transfer_out" => Ok(EventType::TransferOut),
            "transfer_in" => Ok(EventType::TransferIn),
            "adjust" => Ok(EventType::Adjust),
            "earn" => Ok(EventType::Earn),
            "redeem" => Ok(EventType::Redeem),
            _ => Err(format!("Invalid event type: {}", s)),
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, ToSchema)]
pub struct PointLedger {
    pub id: u32,
    #[serde(rename = "userId")]
    pub user_id: u32,
    pub change: i32,
    #[serde(rename = "balanceAfter")]
    pub balance_after: u32,
    #[serde(rename = "eventType")]
    pub event_type: EventType,
    #[serde(rename = "transferId")]
    pub transfer_id: Option<u32>,
    pub reference: Option<String>,
    pub metadata: Option<String>,
    #[serde(rename = "createdAt")]
    #[schema(value_type = String, format = "date-time")]
    pub created_at: DateTime<Utc>,
}

// Database model for internal use
#[derive(Debug, Clone)]
pub struct PointLedgerDb {
    pub id: u32,
    pub user_id: u32,
    pub change: i32,
    pub balance_after: u32,
    pub event_type: String,
    pub transfer_id: Option<u32>,
    pub reference: Option<String>,
    pub metadata: Option<String>,
    pub created_at: String,
}

impl PointLedgerDb {
    pub fn to_domain(self) -> Result<PointLedger, String> {
        let event_type = self.event_type.parse::<EventType>()
            .map_err(|e| format!("Invalid event type: {}", e))?;
        
        let created_at = DateTime::parse_from_rfc3339(&self.created_at)
            .map_err(|e| format!("Invalid created_at date: {}", e))?
            .with_timezone(&Utc);
        
        Ok(PointLedger {
            id: self.id,
            user_id: self.user_id,
            change: self.change,
            balance_after: self.balance_after,
            event_type,
            transfer_id: self.transfer_id,
            reference: self.reference,
            metadata: self.metadata,
            created_at,
        })
    }
}