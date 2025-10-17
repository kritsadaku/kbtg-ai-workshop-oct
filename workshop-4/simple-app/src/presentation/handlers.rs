use axum::{
    extract::{Path, State, Query},
    http::StatusCode,
    response::Json,
};
use serde::{Deserialize, Serialize};
use utoipa::ToSchema;
use crate::application::{UserService, TransferService};
use crate::domain::{User, CreateUserRequest, UpdateUserRequest};

#[derive(Clone)]
pub struct AppState {
    pub user_service: UserService,
    pub transfer_service: TransferService,
}

#[derive(Deserialize)]
pub struct ListUsersQuery {
    pub limit: Option<i64>,
    pub offset: Option<i64>,
}

#[derive(Serialize, ToSchema)]
pub struct ErrorResponse {
    pub error: String,
    pub message: String,
}

#[derive(Serialize, ToSchema)]
pub struct ListUsersResponse {
    pub users: Vec<User>,
    pub total: usize,
}

/// Hello World endpoint
#[utoipa::path(
    get,
    path = "/",
    responses(
        (status = 200, description = "Hello World message", body = String)
    )
)]
pub async fn hello_world() -> &'static str {
    "Hello, World! - Clean Architecture API with SQLite"
}

/// Get user by ID
#[utoipa::path(
    get,
    path = "/users/{id}",
    params(
        ("id" = u32, Path, description = "User ID")
    ),
    responses(
        (status = 200, description = "User found successfully", body = User),
        (status = 404, description = "User not found", body = ErrorResponse)
    )
)]
pub async fn get_user(
    State(state): State<AppState>,
    Path(id): Path<u32>,
) -> Result<Json<User>, (StatusCode, Json<ErrorResponse>)> {
    match state.user_service.get_user(id).await {
        Ok(Some(user)) => Ok(Json(user)),
        Ok(None) => Err((
            StatusCode::NOT_FOUND,
            Json(ErrorResponse {
                error: "USER_NOT_FOUND".to_string(),
                message: "User not found".to_string(),
            }),
        )),
        Err(err) => Err((
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(ErrorResponse { 
                error: "INTERNAL_ERROR".to_string(),
                message: err,
            }),
        )),
    }
}

/// List all users
#[utoipa::path(
    get,
    path = "/users",
    params(
        ("limit" = Option<i64>, Query, description = "Number of users to return (default: 100)"),
        ("offset" = Option<i64>, Query, description = "Number of users to skip (default: 0)")
    ),
    responses(
        (status = 200, description = "List of users", body = ListUsersResponse)
    )
)]
pub async fn list_users(
    State(state): State<AppState>,
    Query(params): Query<ListUsersQuery>,
) -> Result<Json<ListUsersResponse>, (StatusCode, Json<ErrorResponse>)> {
    match state.user_service.list_users(params.limit, params.offset).await {
        Ok(users) => {
            let total = users.len();
            Ok(Json(ListUsersResponse { users, total }))
        }
        Err(err) => Err((
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(ErrorResponse { 
                error: "INTERNAL_ERROR".to_string(),
                message: err,
            }),
        )),
    }
}

/// Create a new user
#[utoipa::path(
    post,
    path = "/users",
    request_body = CreateUserRequest,
    responses(
        (status = 201, description = "User created successfully", body = User),
        (status = 400, description = "Bad request", body = ErrorResponse)
    )
)]
pub async fn create_user(
    State(state): State<AppState>,
    Json(payload): Json<CreateUserRequest>,
) -> Result<(StatusCode, Json<User>), (StatusCode, Json<ErrorResponse>)> {
    match state.user_service.create_user(payload).await {
        Ok(user) => Ok((StatusCode::CREATED, Json(user))),
        Err(err) => Err((
            StatusCode::BAD_REQUEST,
            Json(ErrorResponse { 
                error: "VALIDATION_ERROR".to_string(),
                message: err,
            }),
        )),
    }
}

/// Update an existing user
#[utoipa::path(
    put,
    path = "/users/{id}",
    params(
        ("id" = u32, Path, description = "User ID")
    ),
    request_body = UpdateUserRequest,
    responses(
        (status = 200, description = "User updated successfully", body = User),
        (status = 400, description = "Bad request", body = ErrorResponse),
        (status = 404, description = "User not found", body = ErrorResponse)
    )
)]
pub async fn update_user(
    State(state): State<AppState>,
    Path(id): Path<u32>,
    Json(payload): Json<UpdateUserRequest>,
) -> Result<Json<User>, (StatusCode, Json<ErrorResponse>)> {
    match state.user_service.update_user(id, payload).await {
        Ok(user) => Ok(Json(user)),
        Err(err) => {
            let status = if err.contains("not found") {
                StatusCode::NOT_FOUND
            } else {
                StatusCode::BAD_REQUEST
            };
            Err((status, Json(ErrorResponse { 
                error: if err.contains("not found") { "USER_NOT_FOUND" } else { "VALIDATION_ERROR" }.to_string(),
                message: err,
            })))
        }
    }
}

/// Delete a user
#[utoipa::path(
    delete,
    path = "/users/{id}",
    params(
        ("id" = u32, Path, description = "User ID")
    ),
    responses(
        (status = 204, description = "User deleted successfully"),
        (status = 404, description = "User not found", body = ErrorResponse)
    )
)]
pub async fn delete_user(
    State(state): State<AppState>,
    Path(id): Path<u32>,
) -> Result<StatusCode, (StatusCode, Json<ErrorResponse>)> {
    match state.user_service.delete_user(id).await {
        Ok(true) => Ok(StatusCode::NO_CONTENT),
        Ok(false) => Err((
            StatusCode::NOT_FOUND,
            Json(ErrorResponse {
                error: "USER_NOT_FOUND".to_string(),
                message: "User not found".to_string(),
            }),
        )),
        Err(err) => Err((
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(ErrorResponse { 
                error: "INTERNAL_ERROR".to_string(),
                message: err,
            }),
        )),
    }
}