use axum::{
    extract::{Path, Query, State},
    http::StatusCode,
    response::Json,
};
use serde::Deserialize;
use crate::domain::{CreateTransferRequest, TransferCreateResponse, TransferGetResponse, TransferListResponse};
use crate::presentation::{AppState, ErrorResponse};

#[derive(Deserialize)]
pub struct ListTransfersQuery {
    #[serde(rename = "userId")]
    pub user_id: u32,
    pub page: Option<u32>,
    #[serde(rename = "pageSize")]
    pub page_size: Option<u32>,
}

/// Create a new transfer
#[utoipa::path(
    post,
    path = "/transfers",
    request_body = CreateTransferRequest,
    responses(
        (status = 201, description = "Transfer created successfully", body = TransferCreateResponse),
        (status = 400, description = "Bad request", body = ErrorResponse),
        (status = 409, description = "Conflict (insufficient points)", body = ErrorResponse),
        (status = 422, description = "Unprocessable entity", body = ErrorResponse)
    ),
    tag = "Transfers"
)]
pub async fn create_transfer(
    State(state): State<AppState>,
    Json(request): Json<CreateTransferRequest>,
) -> Result<(StatusCode, Json<TransferCreateResponse>), (StatusCode, Json<ErrorResponse>)> {
    match state.transfer_service.create_transfer(request).await {
        Ok(response) => Ok((StatusCode::CREATED, Json(response))),
        Err(e) => {
            if e.contains("not found") {
                Err((
                    StatusCode::BAD_REQUEST,
                    Json(ErrorResponse {
                        error: "VALIDATION_ERROR".to_string(),
                        message: e,
                    }),
                ))
            } else if e.contains("Insufficient points") {
                Err((
                    StatusCode::CONFLICT,
                    Json(ErrorResponse {
                        error: "INSUFFICIENT_POINTS".to_string(),
                        message: e,
                    }),
                ))
            } else if e.contains("Cannot transfer to the same user") {
                Err((
                    StatusCode::UNPROCESSABLE_ENTITY,
                    Json(ErrorResponse {
                        error: "INVALID_TRANSFER".to_string(),
                        message: e,
                    }),
                ))
            } else {
                Err((
                    StatusCode::BAD_REQUEST,
                    Json(ErrorResponse {
                        error: "VALIDATION_ERROR".to_string(),
                        message: e,
                    }),
                ))
            }
        }
    }
}

/// Get transfer by idempotency key
#[utoipa::path(
    get,
    path = "/transfers/{id}",
    params(
        ("id" = String, Path, description = "Transfer idempotency key")
    ),
    responses(
        (status = 200, description = "Transfer found", body = TransferGetResponse),
        (status = 404, description = "Transfer not found", body = ErrorResponse)
    ),
    tag = "Transfers"
)]
pub async fn get_transfer(
    State(state): State<AppState>,
    Path(id): Path<String>,
) -> Result<Json<TransferGetResponse>, (StatusCode, Json<ErrorResponse>)> {
    match state.transfer_service.get_transfer(&id).await {
        Ok(response) => Ok(Json(response)),
        Err(e) => {
            if e.contains("not found") {
                Err((
                    StatusCode::NOT_FOUND,
                    Json(ErrorResponse {
                        error: "TRANSFER_NOT_FOUND".to_string(),
                        message: e,
                    }),
                ))
            } else {
                Err((
                    StatusCode::BAD_REQUEST,
                    Json(ErrorResponse {
                        error: "VALIDATION_ERROR".to_string(),
                        message: e,
                    }),
                ))
            }
        }
    }
}

/// List transfers for a user
#[utoipa::path(
    get,
    path = "/transfers",
    params(
        ("userId" = u32, Query, description = "User ID to filter transfers"),
        ("page" = Option<u32>, Query, description = "Page number (default: 1)"),
        ("pageSize" = Option<u32>, Query, description = "Page size (default: 20, max: 200)")
    ),
    responses(
        (status = 200, description = "Transfers found", body = TransferListResponse),
        (status = 400, description = "Bad request", body = ErrorResponse)
    ),
    tag = "Transfers"
)]
pub async fn list_transfers(
    State(state): State<AppState>,
    Query(params): Query<ListTransfersQuery>,
) -> Result<Json<TransferListResponse>, (StatusCode, Json<ErrorResponse>)> {
    let page = params.page.unwrap_or(1);
    let page_size = params.page_size.unwrap_or(20);

    match state.transfer_service.list_transfers(params.user_id, page, page_size).await {
        Ok(response) => Ok(Json(response)),
        Err(e) => {
            if e.contains("not found") {
                Err((
                    StatusCode::BAD_REQUEST,
                    Json(ErrorResponse {
                        error: "USER_NOT_FOUND".to_string(),
                        message: e,
                    }),
                ))
            } else {
                Err((
                    StatusCode::BAD_REQUEST,
                    Json(ErrorResponse {
                        error: "VALIDATION_ERROR".to_string(),
                        message: e,
                    }),
                ))
            }
        }
    }
}