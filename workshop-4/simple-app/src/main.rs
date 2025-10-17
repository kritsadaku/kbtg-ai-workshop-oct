mod domain;
mod infrastructure;
mod application;
mod presentation;

use std::sync::Arc;
use utoipa::OpenApi;
use utoipa_swagger_ui::SwaggerUi;
use sqlx::SqlitePool;

use domain::{User, CreateUserRequest, UpdateUserRequest};
use infrastructure::SqliteUserRepository;
use application::UserService;
use presentation::{create_routes, AppState, ErrorResponse, ListUsersResponse};

#[derive(OpenApi)]
#[openapi(
    paths(
        presentation::handlers::hello_world,
        presentation::handlers::get_user,
        presentation::handlers::list_users,
        presentation::handlers::create_user,
        presentation::handlers::update_user,
        presentation::handlers::delete_user,
    ),
    components(
        schemas(User, CreateUserRequest, UpdateUserRequest, ErrorResponse, ListUsersResponse)
    ),
    tags(
        (name = "simple-app", description = "Clean Architecture API with User Management and SQLite")
    )
)]
struct ApiDoc;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Database setup
    let database_url = "sqlite:users.db";
    
    // Create database file if it doesn't exist
    if !std::path::Path::new("users.db").exists() {
        println!("📂 Creating SQLite database file: users.db");
        std::fs::File::create("users.db")?;
    }
    
    let pool = SqlitePool::connect(database_url).await?;
    
    // Infrastructure layer - Repository
    let user_repository = Arc::new(SqliteUserRepository::new(pool));
    
    // Initialize database tables
    user_repository.init_database().await?;
    
    // Application layer - Service
    let user_service = UserService::new(user_repository);
    
    // Application state
    let app_state = AppState { user_service };
    
    // Presentation layer - Routes
    let app = create_routes()
        .merge(SwaggerUi::new("/swagger-ui").url("/api-docs/openapi.json", ApiDoc::openapi()))
        .with_state(app_state);

    println!("🚀 Server running on http://0.0.0.0:3000");
    println!("📚 Swagger UI available at http://0.0.0.0:3000/swagger-ui");
    println!("💾 SQLite database: users.db");
    println!("🔗 API Endpoints:");
    println!("   GET    /");
    println!("   GET    /users?limit=10&offset=0");
    println!("   POST   /users");
    println!("   GET    /users/{{id}}");
    println!("   PUT    /users/{{id}}");
    println!("   DELETE /users/{{id}}");

    // run our app with hyper, listening globally on port 3000
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await?;
    axum::serve(listener, app).await?;
    
    Ok(())
}
