use async_trait::async_trait;
use super::user::{User, CreateUserRequest, UpdateUserRequest};

#[async_trait]
pub trait UserRepository {
    async fn get_user_by_id(&self, id: u32) -> Result<Option<User>, String>;
    async fn get_user_by_email(&self, email: &str) -> Result<Option<User>, String>;
    async fn create_user(&self, user_request: CreateUserRequest) -> Result<User, String>;
    async fn update_user(&self, id: u32, update_request: UpdateUserRequest) -> Result<User, String>;
    async fn delete_user(&self, id: u32) -> Result<bool, String>;
    async fn list_users(&self, limit: Option<i64>, offset: Option<i64>) -> Result<Vec<User>, String>;
}