use async_trait::async_trait;
use sqlx::{SqlitePool, Row};
use chrono::{DateTime, Utc};
use crate::domain::{User, UserRepository, CreateUserRequest, UpdateUserRequest};

#[derive(Clone)]
pub struct SqliteUserRepository {
    pool: SqlitePool,
}

impl SqliteUserRepository {
    pub fn new(pool: SqlitePool) -> Self {
        Self { pool }
    }

    pub async fn init_database(&self) -> Result<(), String> {
        sqlx::query(
            r#"
            CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                first_name TEXT NOT NULL,
                last_name TEXT NOT NULL,
                phone TEXT NOT NULL,
                email TEXT NOT NULL UNIQUE,
                member_since TEXT NOT NULL,
                membership_level TEXT NOT NULL,
                points INTEGER NOT NULL DEFAULT 0,
                created_at TEXT NOT NULL,
                updated_at TEXT NOT NULL
            )
            "#,
        )
        .execute(&self.pool)
        .await
        .map_err(|e| format!("Failed to create users table: {}", e))?;

        // Insert sample data if table is empty
        let count: i64 = sqlx::query_scalar("SELECT COUNT(*) FROM users")
            .fetch_one(&self.pool)
            .await
            .map_err(|e| format!("Failed to count users: {}", e))?;

        if count == 0 {
            self.seed_data().await?;
        }

        Ok(())
    }

    async fn seed_data(&self) -> Result<(), String> {
        let users = vec![
            (
                "John",
                "Doe",
                "+66812345678",
                "john.doe@example.com",
                "Gold",
                1500,
            ),
            (
                "Jane",
                "Smith",
                "+66887654321",
                "jane.smith@example.com",
                "Silver",
                750,
            ),
            (
                "Bob",
                "Johnson",
                "+66856789012",
                "bob.johnson@example.com",
                "Bronze",
                200,
            ),
        ];

        for (first_name, last_name, phone, email, membership_level, points) in users {
            let now = Utc::now().to_rfc3339();
            sqlx::query(
                r#"
                INSERT INTO users (first_name, last_name, phone, email, member_since, membership_level, points, created_at, updated_at)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
                "#,
            )
            .bind(first_name)
            .bind(last_name)
            .bind(phone)
            .bind(email)
            .bind(&now)
            .bind(membership_level)
            .bind(points)
            .bind(&now)
            .bind(&now)
            .execute(&self.pool)
            .await
            .map_err(|e| format!("Failed to insert seed data: {}", e))?;
        }

        Ok(())
    }
}

#[async_trait]
impl UserRepository for SqliteUserRepository {
    async fn get_user_by_id(&self, id: u32) -> Result<Option<User>, String> {
        let row = sqlx::query(
            "SELECT id, first_name, last_name, phone, email, member_since, membership_level, points, created_at, updated_at FROM users WHERE id = ?"
        )
        .bind(id as i64)
        .fetch_optional(&self.pool)
        .await
        .map_err(|e| format!("Database error: {}", e))?;

        match row {
            Some(row) => {
                let user = User {
                    id: row.get::<i64, _>("id") as u32,
                    first_name: row.get("first_name"),
                    last_name: row.get("last_name"),
                    phone: row.get("phone"),
                    email: row.get("email"),
                    member_since: DateTime::parse_from_rfc3339(row.get("member_since"))
                        .map_err(|e| format!("Invalid datetime: {}", e))?
                        .with_timezone(&Utc),
                    membership_level: row.get("membership_level"),
                    points: row.get("points"),
                    created_at: DateTime::parse_from_rfc3339(row.get("created_at"))
                        .map_err(|e| format!("Invalid datetime: {}", e))?
                        .with_timezone(&Utc),
                    updated_at: DateTime::parse_from_rfc3339(row.get("updated_at"))
                        .map_err(|e| format!("Invalid datetime: {}", e))?
                        .with_timezone(&Utc),
                };
                Ok(Some(user))
            }
            None => Ok(None),
        }
    }

    async fn get_user_by_email(&self, email: &str) -> Result<Option<User>, String> {
        let row = sqlx::query(
            "SELECT id, first_name, last_name, phone, email, member_since, membership_level, points, created_at, updated_at FROM users WHERE email = ?"
        )
        .bind(email)
        .fetch_optional(&self.pool)
        .await
        .map_err(|e| format!("Database error: {}", e))?;

        match row {
            Some(row) => {
                let user = User {
                    id: row.get::<i64, _>("id") as u32,
                    first_name: row.get("first_name"),
                    last_name: row.get("last_name"),
                    phone: row.get("phone"),
                    email: row.get("email"),
                    member_since: DateTime::parse_from_rfc3339(row.get("member_since"))
                        .map_err(|e| format!("Invalid datetime: {}", e))?
                        .with_timezone(&Utc),
                    membership_level: row.get("membership_level"),
                    points: row.get("points"),
                    created_at: DateTime::parse_from_rfc3339(row.get("created_at"))
                        .map_err(|e| format!("Invalid datetime: {}", e))?
                        .with_timezone(&Utc),
                    updated_at: DateTime::parse_from_rfc3339(row.get("updated_at"))
                        .map_err(|e| format!("Invalid datetime: {}", e))?
                        .with_timezone(&Utc),
                };
                Ok(Some(user))
            }
            None => Ok(None),
        }
    }

    async fn create_user(&self, user_request: CreateUserRequest) -> Result<User, String> {
        // Check if email already exists
        if let Some(_) = self.get_user_by_email(&user_request.email).await? {
            return Err("Email already exists".to_string());
        }

        let user = User::new(
            0, // Will be replaced by auto-increment
            user_request.first_name,
            user_request.last_name,
            user_request.phone,
            user_request.email,
            user_request.membership_level,
        );

        user.validate()?;

        let result = sqlx::query(
            r#"
            INSERT INTO users (first_name, last_name, phone, email, member_since, membership_level, points, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
            "#,
        )
        .bind(&user.first_name)
        .bind(&user.last_name)
        .bind(&user.phone)
        .bind(&user.email)
        .bind(user.member_since.to_rfc3339())
        .bind(&user.membership_level)
        .bind(user.points)
        .bind(user.created_at.to_rfc3339())
        .bind(user.updated_at.to_rfc3339())
        .execute(&self.pool)
        .await
        .map_err(|e| format!("Failed to create user: {}", e))?;

        let id = result.last_insert_rowid() as u32;
        
        // Return user with the generated ID
        let mut created_user = user;
        created_user.id = id;
        Ok(created_user)
    }

    async fn update_user(&self, id: u32, update_request: UpdateUserRequest) -> Result<User, String> {
        let mut user = self
            .get_user_by_id(id)
            .await?
            .ok_or("User not found".to_string())?;

        // Check if email is being updated and if it already exists
        if let Some(ref new_email) = update_request.email {
            if new_email != &user.email {
                if let Some(_) = self.get_user_by_email(new_email).await? {
                    return Err("Email already exists".to_string());
                }
            }
        }

        user.update_fields(update_request);
        user.validate()?;

        sqlx::query(
            r#"
            UPDATE users 
            SET first_name = ?, last_name = ?, phone = ?, email = ?, membership_level = ?, updated_at = ?
            WHERE id = ?
            "#,
        )
        .bind(&user.first_name)
        .bind(&user.last_name)
        .bind(&user.phone)
        .bind(&user.email)
        .bind(&user.membership_level)
        .bind(user.updated_at.to_rfc3339())
        .bind(id as i64)
        .execute(&self.pool)
        .await
        .map_err(|e| format!("Failed to update user: {}", e))?;

        Ok(user)
    }

    async fn delete_user(&self, id: u32) -> Result<bool, String> {
        let result = sqlx::query("DELETE FROM users WHERE id = ?")
            .bind(id as i64)
            .execute(&self.pool)
            .await
            .map_err(|e| format!("Failed to delete user: {}", e))?;

        Ok(result.rows_affected() > 0)
    }

    async fn list_users(&self, limit: Option<i64>, offset: Option<i64>) -> Result<Vec<User>, String> {
        let limit = limit.unwrap_or(100);
        let offset = offset.unwrap_or(0);

        let rows = sqlx::query(
            "SELECT id, first_name, last_name, phone, email, member_since, membership_level, points, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?"
        )
        .bind(limit)
        .bind(offset)
        .fetch_all(&self.pool)
        .await
        .map_err(|e| format!("Database error: {}", e))?;

        let mut users = Vec::new();
        for row in rows {
            let user = User {
                id: row.get::<i64, _>("id") as u32,
                first_name: row.get("first_name"),
                last_name: row.get("last_name"),
                phone: row.get("phone"),
                email: row.get("email"),
                member_since: DateTime::parse_from_rfc3339(row.get("member_since"))
                    .map_err(|e| format!("Invalid datetime: {}", e))?
                    .with_timezone(&Utc),
                membership_level: row.get("membership_level"),
                points: row.get("points"),
                created_at: DateTime::parse_from_rfc3339(row.get("created_at"))
                    .map_err(|e| format!("Invalid datetime: {}", e))?
                    .with_timezone(&Utc),
                updated_at: DateTime::parse_from_rfc3339(row.get("updated_at"))
                    .map_err(|e| format!("Invalid datetime: {}", e))?
                    .with_timezone(&Utc),
            };
            users.push(user);
        }

        Ok(users)
    }
}