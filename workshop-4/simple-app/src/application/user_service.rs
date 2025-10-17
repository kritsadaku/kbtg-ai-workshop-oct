use std::sync::Arc;
use crate::domain::{User, UserRepository, CreateUserRequest, UpdateUserRequest};

#[derive(Clone)]
pub struct UserService {
    repository: Arc<dyn UserRepository + Send + Sync>,
}

impl UserService {
    pub fn new(repository: Arc<dyn UserRepository + Send + Sync>) -> Self {
        Self { repository }
    }

    pub async fn get_user(&self, id: u32) -> Result<Option<User>, String> {
        self.repository.get_user_by_id(id).await
    }

    pub async fn get_user_by_email(&self, email: &str) -> Result<Option<User>, String> {
        self.repository.get_user_by_email(email).await
    }

    pub async fn create_user(&self, user_request: CreateUserRequest) -> Result<User, String> {
        self.repository.create_user(user_request).await
    }

    pub async fn update_user(&self, id: u32, update_request: UpdateUserRequest) -> Result<User, String> {
        self.repository.update_user(id, update_request).await
    }

    pub async fn delete_user(&self, id: u32) -> Result<bool, String> {
        self.repository.delete_user(id).await
    }

    pub async fn list_users(&self, limit: Option<i64>, offset: Option<i64>) -> Result<Vec<User>, String> {
        self.repository.list_users(limit, offset).await
    }
}