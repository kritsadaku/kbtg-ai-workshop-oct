use serde::{Deserialize, Serialize};
use utoipa::ToSchema;
use chrono::{DateTime, Utc};

#[derive(Debug, Clone, Serialize, Deserialize, ToSchema)]
pub struct User {
    pub id: u32,
    pub first_name: String,
    pub last_name: String,
    pub phone: String,
    pub email: String, // Unique
    #[serde(with = "chrono::serde::ts_seconds")]
    #[schema(value_type = String, format = DateTime)]
    pub member_since: DateTime<Utc>,
    pub membership_level: String,
    pub points: i64,
    #[serde(with = "chrono::serde::ts_seconds")]
    #[schema(value_type = String, format = DateTime)]
    pub created_at: DateTime<Utc>,
    #[serde(with = "chrono::serde::ts_seconds")]
    #[schema(value_type = String, format = DateTime)]
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize, ToSchema)]
pub struct CreateUserRequest {
    pub first_name: String,
    pub last_name: String,
    pub phone: String,
    pub email: String,
    pub membership_level: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize, ToSchema)]
pub struct UpdateUserRequest {
    pub first_name: Option<String>,
    pub last_name: Option<String>,
    pub phone: Option<String>,
    pub email: Option<String>,
    pub membership_level: Option<String>,
}

impl User {
    pub fn new(
        id: u32,
        first_name: String,
        last_name: String,
        phone: String,
        email: String,
        membership_level: Option<String>,
    ) -> Self {
        let now = Utc::now();
        Self {
            id,
            first_name,
            last_name,
            phone,
            email,
            member_since: now,
            membership_level: membership_level.unwrap_or_else(|| "Bronze".to_string()),
            points: 0,
            created_at: now,
            updated_at: now,
        }
    }

    pub fn validate(&self) -> Result<(), String> {
        if self.first_name.trim().is_empty() {
            return Err("First name cannot be empty".to_string());
        }
        if self.last_name.trim().is_empty() {
            return Err("Last name cannot be empty".to_string());
        }
        if self.phone.trim().is_empty() {
            return Err("Phone cannot be empty".to_string());
        }
        if self.email.trim().is_empty() {
            return Err("Email cannot be empty".to_string());
        }
        if !self.email.contains('@') {
            return Err("Invalid email format".to_string());
        }
        Ok(())
    }

    pub fn update_fields(&mut self, update_request: UpdateUserRequest) {
        if let Some(first_name) = update_request.first_name {
            self.first_name = first_name;
        }
        if let Some(last_name) = update_request.last_name {
            self.last_name = last_name;
        }
        if let Some(phone) = update_request.phone {
            self.phone = phone;
        }
        if let Some(email) = update_request.email {
            self.email = email;
        }
        if let Some(membership_level) = update_request.membership_level {
            self.membership_level = membership_level;
        }
        self.updated_at = Utc::now();
    }
}