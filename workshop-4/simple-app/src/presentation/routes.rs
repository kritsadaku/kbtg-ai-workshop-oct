use axum::{
    routing::{get, post, put, delete},
    Router,
};
use super::handlers::{
    hello_world, get_user, list_users, create_user, update_user, delete_user, AppState
};
use super::transfer_handlers::{
    create_transfer, get_transfer, list_transfers
};

pub fn create_routes() -> Router<AppState> {
    Router::new()
        .route("/", get(hello_world))
        .route("/users", get(list_users))
        .route("/users", post(create_user))
        .route("/users/{id}", get(get_user))
        .route("/users/{id}", put(update_user))
        .route("/users/{id}", delete(delete_user))
        .route("/transfers", post(create_transfer))
        .route("/transfers", get(list_transfers))
        .route("/transfers/{id}", get(get_transfer))
}