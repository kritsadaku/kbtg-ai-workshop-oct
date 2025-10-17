use axum::{
    routing::{get, post, put, delete},
    Router,
};
use super::handlers::{
    hello_world, get_user, list_users, create_user, update_user, delete_user, AppState
};

pub fn create_routes() -> Router<AppState> {
    Router::new()
        .route("/", get(hello_world))
        .route("/users", get(list_users))
        .route("/users", post(create_user))
        .route("/users/{id}", get(get_user))
        .route("/users/{id}", put(update_user))
        .route("/users/{id}", delete(delete_user))
}